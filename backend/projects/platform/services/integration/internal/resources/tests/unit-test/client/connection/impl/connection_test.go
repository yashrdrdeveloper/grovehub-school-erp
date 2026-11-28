package impl_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"grovehub/services/integration/client/connection"
	clientimpl "grovehub/services/integration/internal/app/client/connection/impl"
	connections "grovehub/services/integration/internal/app/services/connection"
	providers "grovehub/services/integration/internal/app/services/provider"
)

type connectionService struct {
	connections.ConnectionService
	ctx context.Context
	tenant string
	id string
	input connection.CreateConnection
	action connection.Action
	payload json.RawMessage
	err error
}

func (s *connectionService) Create(ctx context.Context, tenant string, input connection.CreateConnection) (connection.Connection, error) {
	s.ctx, s.tenant, s.input = ctx, tenant, input
	return connection.Connection{ID: "created"}, s.err
}

func (s *connectionService) Get(ctx context.Context, tenant, id string) (connection.Connection, error) {
	s.ctx, s.tenant, s.id = ctx, tenant, id
	return connection.Connection{ID: id}, s.err
}

func (s *connectionService) List(ctx context.Context, tenant, owner string) ([]connection.Connection, error) {
	s.ctx, s.tenant, s.id = ctx, tenant, owner
	return []connection.Connection{{ID: owner}}, s.err
}

func (s *connectionService) Disconnect(ctx context.Context, tenant, id string) (connection.Connection, error) {
	return s.Get(ctx, tenant, id)
}

func (s *connectionService) ExecuteAction(ctx context.Context, tenant, id string, action connection.Action, input json.RawMessage) (any, error) {
	s.ctx, s.tenant, s.id, s.action, s.payload = ctx, tenant, id, action, input
	return "action", s.err
}

func (s *connectionService) ExecuteSystemMail(ctx context.Context, tenant string, input json.RawMessage) (any, error) {
	s.ctx, s.tenant, s.payload = ctx, tenant, input
	return "mail", s.err
}

type providerService struct {
	providers.ProviderService
	ctx context.Context
	id string
	err error
}

func (s *providerService) Providers(ctx context.Context) []connection.Descriptor {
	s.ctx = ctx
	return []connection.Descriptor{{ID: "provider"}}
}

func (s *providerService) Provider(ctx context.Context, id string) (connection.ProviderDetails, error) {
	s.ctx, s.id = ctx, id
	return connection.ProviderDetails{}, s.err
}

func TestConnectionClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	payload := json.RawMessage("{}")
	for _, failure := range []error{nil, errors.New("service failed")} {
		connections := &connectionService{err: failure}
		providers := &providerService{err: failure}
		client := clientimpl.NewConnectionClient(connections, providers)
		if got := client.Providers(ctx); len(got) != 1 || got[0].ID != "provider" || providers.ctx != ctx {
			t.Fatal("provider list not forwarded")
		}
		if _, err := client.Provider(ctx, "provider"); err != failure || providers.id != "provider" || providers.ctx != ctx {
			t.Fatal("provider lookup not forwarded")
		}
		if got, err := client.CreateConnection(ctx, "tenant", connection.CreateConnection{}); err != failure || got.ID != "created" || connections.tenant != "tenant" || connections.ctx != ctx {
			t.Fatal("create not forwarded")
		}
		if got, err := client.GetConnection(ctx, "tenant", "id"); err != failure || got.ID != "id" || connections.id != "id" {
			t.Fatal("get not forwarded")
		}
		if got, err := client.ListConnections(ctx, "tenant", "owner"); err != failure || len(got) != 1 || got[0].ID != "owner" {
			t.Fatal("list not forwarded")
		}
		if got, err := client.DisconnectConnection(ctx, "tenant", "id"); err != failure || got.ID != "id" {
			t.Fatal("disconnect not forwarded")
		}
		if got, err := client.ExecuteAction(ctx, "tenant", "id", connection.Action("mail.send"), payload); err != failure || got != "action" || connections.action != "mail.send" || !reflect.DeepEqual(connections.payload, payload) {
			t.Fatal("action not forwarded")
		}
		if got, err := client.ExecuteSystemMail(ctx, "tenant", payload); err != failure || got != "mail" || !reflect.DeepEqual(connections.payload, payload) {
			t.Fatal("system mail not forwarded")
		}
	}
}
