package audit

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
	sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/audit"
	core "grovehub/libs/core/gateway"
	"net/http"
)

type ReceiptSearchInput struct {
	entity.PageRequest
	SchoolID  string `json:"schoolId,omitempty"`
	Operation string `json:"operation,omitempty"`
}

func (in ReceiptSearchInput) Normalize() (ReceiptSearchInput, bool) {
	page, ok := in.PageRequest.Normalize()
	in.PageRequest = page
	if !ok {
		return in, false
	}
	if in.SchoolID != "" {
		id, err := uuid.Parse(in.SchoolID)
		if err != nil || id == uuid.Nil || id.String() != in.SchoolID {
			return in, false
		}
	}
	switch in.Operation {
	case "", "RENAME", "ACTIVATE", "DISABLE":
		return in, true
	default:
		return in, false
	}
}

type ReceiptReader func(context.Context, ReceiptSearchInput) (*sdk.SchoolReceiptPage, error)

func NewReceiptReader(c Config) (ReceiptReader, error) {
	api, err := buildClient(c, "internal-tool-school-receipt-reader", receiptReadPolicy)
	if err != nil || api == nil {
		return nil, err
	}
	return func(ctx context.Context, p ReceiptSearchInput) (*sdk.SchoolReceiptPage, error) {
		p, ok := p.Normalize()
		if !ok {
			return nil, errors.New("invalid receipt search")
		}
		input := sdk.SchoolReceiptSearchInput{PageNumber: int32(p.PageNumber), RowSize: int32(p.RowSize)}
		if p.SchoolID != "" {
			input.SetSchoolId(p.SchoolID)
		}
		if p.Operation != "" {
			input.SetOperation(p.Operation)
		}
		page, response, err := api.DefaultAPI.SearchSchoolMutationReceipts(ctx).SchoolReceiptSearchInput(input).Execute()
		if response != nil && response.Body != nil {
			defer response.Body.Close()
		}
		if err != nil || response == nil || response.StatusCode != 200 || page == nil {
			return nil, errors.New("central audit unavailable")
		}
		return page, nil
	}, nil
}
func receiptReadPolicy(r *http.Request) (core.WorkloadClaims, error) {
	if r == nil || r.URL == nil || r.Method != "POST" || r.URL.Path != "/audit/restricted/school-mutation-receipts/search" || r.URL.RawPath != "" || r.URL.RawQuery != "" || r.URL.ForceQuery || r.Header.Get("Cookie") != "" || r.Header.Get("X-GH-Identity-Proof") != "" {
		return core.WorkloadClaims{}, core.ErrWorkloadDenied
	}
	return core.WorkloadClaims{Scopes: []string{"internal.audit.school.mutation-receipt.read"}}, nil
}
