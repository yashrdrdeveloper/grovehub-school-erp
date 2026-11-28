package helpers

import (
	"context"

	"gorm.io/gorm"
)

type contextKey struct{}

type State struct {
	db          *gorm.DB
	afterCommit []func()
}

func FromContext(ctx context.Context) (*State, bool) {
	state, ok := ctx.Value(contextKey{}).(*State)

	return state, ok
}

func Bind(ctx context.Context, db *gorm.DB) (context.Context, *State) {
	state := &State{db: db}

	return context.WithValue(ctx, contextKey{}, state), state
}

func (state *State) Database(ctx context.Context) *gorm.DB {
	return state.db.WithContext(ctx)
}

func (state *State) AddCallback(callback func()) {
	state.afterCommit = append(state.afterCommit, callback)
}

func (state *State) PublishCallbacks() {
	for _, callback := range state.afterCommit {
		callback()
	}
}

func DetachedContext(ctx context.Context) context.Context {
	return context.WithValue(context.WithoutCancel(ctx), contextKey{}, nil)
}
