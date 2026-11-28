package impl

import (
	"context"
	"encoding/json"
	"fmt"
	contracts "grovehub/services/integration/internal/app/providers/contracts"
    dto "grovehub/services/integration/internal/app/providers/dto"
    helpers "grovehub/services/integration/internal/app/providers/helpers"
	"sort"
	"strings"
)

type actionRegistry struct {
	handlers map[string]map[dto.Action]contracts.ActionHandler
}

func NewActionRegistry() contracts.ActionRegistry {
	return &actionRegistry{handlers: make(map[string]map[dto.Action]contracts.ActionHandler)}
}

func (r *actionRegistry) Register(provider string, action dto.Action, handler contracts.ActionHandler) error {
	provider = strings.ToUpper(strings.TrimSpace(provider))
	if provider == "" || strings.TrimSpace(string(action)) == "" || handler == nil {
		return helpers.ErrInvalidActionInput
	}
	if r.handlers[provider] == nil {
		r.handlers[provider] = make(map[dto.Action]contracts.ActionHandler)
	}
	if _, exists := r.handlers[provider][action]; exists {
		return fmt.Errorf("%w: %s %s", helpers.ErrDuplicateProvider, provider, action)
	}
	r.handlers[provider][action] = handler
	return nil
}

func (r *actionRegistry) Execute(ctx context.Context, session dto.Session, action dto.Action, input json.RawMessage) (any, error) {
	items := r.handlers[strings.ToUpper(session.ProviderID)]
	if items == nil {
		return nil, helpers.ErrProviderUnavailable
	}
	handler, ok := items[action]
	if !ok {
		return nil, helpers.ErrUnsupportedAction
	}
	return handler.Execute(ctx, session, input)
}

func (r *actionRegistry) Actions(provider string) []dto.Action {
	handlers := r.handlers[strings.ToUpper(provider)]
	result := make([]dto.Action, 0, len(handlers))
	for action := range handlers {
		result = append(result, action)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}
