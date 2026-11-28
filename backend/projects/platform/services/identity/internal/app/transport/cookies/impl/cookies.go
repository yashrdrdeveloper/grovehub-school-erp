// Package cookies owns HTTP cookie encoding only. Authentication, rate limits,
// Origin checks, and CSRF enforcement are implemented in Launcher gateway.
package impl

import (
	"net/http"
	"strings"

	cookies "grovehub/services/identity/internal/app/transport/cookies"
)

const (
	AccessCookie = cookies.AccessCookie
	RefreshCookie = cookies.RefreshCookie
	CSRFName = cookies.CSRFName
	ResetCookie = cookies.ResetCookie
	SetupCookie = cookies.SetupCookie
)

type Manager struct{}

var _ cookies.Manager = (*Manager)(nil)

func New() cookies.Manager {
	return &Manager{}
}


func (*Manager) Issue(writer http.ResponseWriter, access string, refresh string, csrf string) {
	http.SetCookie(writer, &http.Cookie{
		Name: AccessCookie,
		Value: access,
		Path: cookies.SessionPath,
		MaxAge: 900,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(writer, &http.Cookie{
		Name: RefreshCookie,
		Value: refresh,
		Path: cookies.RefreshPath,
		MaxAge: 30 * 24 * 60 * 60,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	if csrf != "" {
		http.SetCookie(writer, &http.Cookie{
			Name: CSRFName,
			Value: csrf,
			Path: cookies.SessionPath,
			MaxAge: 30 * 24 * 60 * 60,
			Secure: true,
			SameSite: http.SameSiteStrictMode,
		})
	}
	writer.Header().Set(cookies.HeaderCacheControl, cookies.DirectiveNoStore)
}

func (*Manager) ClearSession(writer http.ResponseWriter) {
	paths := []struct {
		name string
		path string
	}{
		{AccessCookie, cookies.SessionPath},
		{RefreshCookie, cookies.RefreshPath},
		{CSRFName, cookies.SessionPath},
	}
	for _, cookie := range paths {
		http.SetCookie(writer, &http.Cookie{
			Name: cookie.name,
			Path: cookie.path,
			MaxAge: -1,
			Secure: true,
			HttpOnly: cookie.name != CSRFName,
			SameSite: http.SameSiteStrictMode,
		})
	}
	writer.Header().Set(cookies.HeaderCacheControl, cookies.DirectiveNoStore)
}

func (*Manager) Read(request *http.Request, name string) string {
	cookie, err := request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (*Manager) SetGrant(writer http.ResponseWriter, name string, path string, value string) {
	http.SetCookie(writer, &http.Cookie{
		Name: name,
		Path: path,
		Value: value,
		MaxAge: 600,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	writer.Header().Set(cookies.HeaderCacheControl, cookies.DirectiveNoStore)
}

func (m *Manager) ClearGrant(writer http.ResponseWriter, name string, path string) {
	http.SetCookie(writer, &http.Cookie{
		Name: name,
		Path: path,
		MaxAge: -1,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	writer.Header().Set(cookies.HeaderCacheControl, cookies.DirectiveNoStore)
}

func (*Manager) Proof(request *http.Request, prefix string) string {
	authorization := request.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, prefix+" ") {
		return ""
	}
	return strings.TrimPrefix(authorization, prefix+" ")
}
