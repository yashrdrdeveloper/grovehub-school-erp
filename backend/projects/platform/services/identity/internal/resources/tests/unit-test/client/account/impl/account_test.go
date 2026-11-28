package impl_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	corehttp "grovehub/libs/core/http"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/identity/client/account"
	clientimpl "grovehub/services/identity/internal/app/client/account/impl"
)

func TestAccountClient(t *testing.T) {
	const id = "account-id"
	client := clientimpl.NewAccountClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if !strings.HasSuffix(r.URL.Path, "/restricted/accounts/"+id) {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if err := corehttp.WriteSuccess(w, http.StatusOK, account.Account{AccountID: id}); err != nil {
				t.Fatal(err)
			}
		case http.MethodPost:
			if !strings.HasSuffix(r.URL.Path, "/restricted/accounts/search") {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if err := corehttp.WriteSuccess(w, http.StatusOK, coresearch.Result[account.Account]{}); err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	}))
	got, err := client.Get(context.Background(), account.GetRequest{Params: account.ByID{AccountID: id}})
	if err != nil || got.AccountID != id {
		t.Fatalf("account=%v error=%v", got, err)
	}
	if _, err := client.Search(context.Background(), account.SearchRequest{}); err != nil {
		t.Fatal(err)
	}
}
