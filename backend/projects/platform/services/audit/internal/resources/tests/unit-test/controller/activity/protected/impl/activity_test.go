package impl_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	controller "grovehub/services/audit/internal/app/controller/activity/protected/impl"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/services/activity"
)

type searchService struct {
	err error
}

func (s searchService) Search(context.Context, dto.Search) (dto.Page, error) {
	return dto.Page{}, s.err
}

func TestSearchErrorContract(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"success", nil, http.StatusOK},
		{"forbidden", activity.ErrForbidden, http.StatusForbidden},
		{"invalid", activity.ErrInvalidSearch, http.StatusBadRequest},
		{"unexpected", errors.New("private database error"), http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/audit/protected/search", strings.NewReader("{}"))
			controller.New(searchService{err: tc.err}).Search(w, r)
			if w.Code != tc.status || strings.Contains(w.Body.String(), "private database error") {
				t.Fatalf("unexpected status=%d body=%s", w.Code, w.Body.String())
			}
			if w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("missing no-store header")
			}
		})
	}
}
