package impl

import (
	"context"
	"errors"

	coresearch "grovehub/libs/core/search"
	dto "grovehub/services/identity/internal/app/dto/approval"
	repo "grovehub/services/identity/internal/app/repository/approval"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/services/approval"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
	validator "grovehub/services/identity/internal/app/validators/approval"
)

type ReadService struct {
	reader repo.Reader
	repository repo.Repository
	locker authority.Locker
	snapshots authrepo.Repository
}

var _ approval.Reader = (*ReadService)(nil)

func NewReader(reader repo.Reader, repository repo.Repository, locker authority.Locker, snapshots authrepo.Repository) approval.Reader {
	return &ReadService{reader: reader, repository: repository, locker: locker, snapshots: snapshots}
}

func (service *ReadService) read(ctx context.Context, operation func(context.Context, repo.Visibility) error) error {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if service.reader == nil || service.repository == nil || service.locker == nil || service.snapshots == nil {
		return authorization.ErrUnavailable
	}
	return service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		snapshot, err := service.snapshots.LoadSnapshot(tx, actor.TenantID, actor.AccountID)
		if err != nil {
			return err
		}
		grants, err := helpers.EligibleGrants(actor.AccountID, actor.TenantID, snapshot)
		if err != nil {
			return err
		}
		now, err := service.repository.Now(tx)
		if err != nil {
			return err
		}
		return operation(tx, repo.Visibility{TenantID: actor.TenantID, AccountID: actor.AccountID, Grants: grants, Now: now})
	})
}

func (service *ReadService) History(ctx context.Context, request coresearch.Request) (coresearch.Result[dto.Summary], error) {
	return service.search(ctx, request, false)
}

func (service *ReadService) Inbox(ctx context.Context, request coresearch.Request) (coresearch.Result[dto.Summary], error) {
	return service.search(ctx, request, true)
}

func (service *ReadService) search(ctx context.Context, request coresearch.Request, inbox bool) (coresearch.Result[dto.Summary], error) {
	query, err := validator.NormalizeSearch(request)
	if err != nil {
		return coresearch.Result[dto.Summary]{}, err
	}
	var result coresearch.Result[dto.Summary]
	err = service.read(ctx, func(tx context.Context, scope repo.Visibility) error {
		scope.Inbox = inbox
		rows, err := service.reader.Search(tx, scope, query)
		if err != nil {
			return err
		}
		result, err = helpers.Hydrate(query, rows)
		return err
	})
	if err != nil {
		return coresearch.Result[dto.Summary]{}, err
	}
	return result, nil
}

func (service *ReadService) PendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := service.read(ctx, func(tx context.Context, scope repo.Visibility) error {
		var err error
		count, err = service.reader.CountInbox(tx, scope)
		return err
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (service *ReadService) Detail(ctx context.Context, requestID string) (dto.Detail, error) {
	id, err := helpers.NormalizeID(requestID)
	if err != nil {
		return dto.Detail{}, err
	}
	var detail dto.Detail
	err = service.read(ctx, func(tx context.Context, scope repo.Visibility) error {
		request, err := service.reader.Get(tx, scope, id)
		if errors.Is(err, repo.ErrNotFound) {
			scope.Inbox = true
			request, err = service.reader.Get(tx, scope, id)
		}
		if err != nil {
			return err
		}
		if err := helpers.VerifyIntent(request); err != nil {
			return err
		}
		if scope.Inbox {
			if err := helpers.ReviewerIndependent(request, scope.AccountID); err != nil {
				return repo.ErrNotFound
			}
		}
		detail = dto.Detail{Summary: helpers.ToSummary(request), PolicyID: request.PolicyID, Intent: request.Payload, Reason: request.Reason}
		return nil
	})
	if err != nil {
		return dto.Detail{}, err
	}
	return detail, nil
}
