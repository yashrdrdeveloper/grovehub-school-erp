package constants_test

import (
	"testing"

	"grovehub/libs/platform/session/internal/src/constants"
)

func TestCookieNames(t *testing.T) {
	got := []string{constants.AccessCookie, constants.RefreshCookie, constants.CSRFName, constants.ResetCookie, constants.SetupCookie}
	want := []string{"__Host-gh_access", "__Secure-gh_refresh", "gh_csrf", "__Secure-gh_reset", "__Secure-gh_setup"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("cookie %d = %q, want %q", i, got[i], want[i])
		}
	}
}
