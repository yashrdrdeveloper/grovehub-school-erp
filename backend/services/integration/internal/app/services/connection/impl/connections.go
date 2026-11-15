package impl

import (
    watchrepository "grovehub/services/integration/internal/app/repository/watch"
	"context"
	"errors"
	"time"

	coreid "grovehub/libs/core/id"
	infrastructureclient "grovehub/services/infrastructure/client/messaging"
	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	
	"grovehub/services/integration/internal/app/entity"
	"grovehub/services/integration/internal/app/providers"
	"grovehub/services/integration/internal/app/providers/mailprovider/providers/gmail"
	"grovehub/services/integration/internal/app/providers/mailprovider/providers/outlook"
	integrationrepository "grovehub/services/integration/internal/app/repository/connection"
	integrationservices "grovehub/services/integration/internal/app/services/connection"
	connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
)

var (
	ErrInvalidConnection  = errors.New("invalid integration connection")
	ErrConnectionNotReady = errors.New("provider authorization is not implemented")
)

type ConnectionService struct {
	repository      integrationrepository.ConnectionRepository
	watchRepository watchrepository.WatchRepository
	pushTopic       string
	messaging       infrastructureclient.MessagingClient
	registry        providers.Registry
	actions         providers.ActionRegistry
	oauth           providers.OAuthRegistry
	gmail           *gmail.Client
	outlook         *outlook.Client
	cipher          tokenCipher
}

func NewConnectionService(
	repository integrationrepository.ConnectionRepository,
	watchRepository watchrepository.WatchRepository,
	registry providers.Registry,
	actions providers.ActionRegistry,
	oauth providers.OAuthRegistry,
	gmailClient *gmail.Client,
	outlookClient *outlook.Client,
	encryptionKey []byte,
	pushTopic string,
	messaging infrastructureclient.MessagingClient,
) (integrationservices.ConnectionService, error) {
	cipher, err := newTokenCipher(encryptionKey)
	if err != nil && (gmailClient != nil || outlookClient != nil) {
		return nil, err
	}

	return &ConnectionService{
		repository:      repository,
		watchRepository: watchRepository,
		pushTopic:       pushTopic,
		messaging:       messaging,
		registry:        registry,
		actions:         actions,
		oauth:           oauth,
		gmail:           gmailClient,
		outlook:         outlookClient,
		cipher:          cipher,
	}, nil
}

func (s *ConnectionService) Create(
	ctx context.Context,
	tenantID string,
	input connectiondto.CreateConnection,
) (connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(input.OwnerGID) {
		return connectiondto.Connection{}, ErrInvalidConnection
	}

	provider, err := s.registry.Find(input.Provider)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	id, err := coreid.NewUUIDv7()
	if err != nil {
		return connectiondto.Connection{}, err
	}

	now := time.Now().UTC()
	connection := entity.Connection{
		ID:        id,
		TenantID:  tenantID,
		OwnerGID:  input.OwnerGID,
		Provider:  provider.Descriptor().ID,
		Status:    entity.StatusPendingSetup,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(ctx, connection); err != nil {
		return connectiondto.Connection{}, err
	}

	return connectionhelpers.ToResponse(connection), nil
}

func (s *ConnectionService) Get(
	ctx context.Context,
	tenantID string,
	id string,
) (connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(id) {
		return connectiondto.Connection{}, ErrInvalidConnection
	}

	connection, err := s.repository.Get(ctx, tenantID, id)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	return connectionhelpers.ToResponse(connection), nil
}

func (s *ConnectionService) List(
	ctx context.Context,
	tenantID string,
	ownerGID string,
) ([]connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(ownerGID) {
		return nil, ErrInvalidConnection
	}

	connections, err := s.repository.List(ctx, tenantID, ownerGID)
	if err != nil {
		return nil, err
	}

	result := make([]connectiondto.Connection, 0, len(connections))
	for _, connection := range connections {
		result = append(result, connectionhelpers.ToResponse(connection))
	}

	return result, nil
}

func (s *ConnectionService) Disconnect(
	ctx context.Context,
	tenantID string,
	id string,
) (connectiondto.Connection, error) {
	if !connectionhelpers.ValidID(tenantID) || !connectionhelpers.ValidID(id) {
		return connectiondto.Connection{}, ErrInvalidConnection
	}

	connection, err := s.repository.Disconnect(ctx, tenantID, id)
	if err != nil {
		return connectiondto.Connection{}, err
	}

	return connectionhelpers.ToResponse(connection), nil
}
