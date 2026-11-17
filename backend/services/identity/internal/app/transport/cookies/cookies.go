// Package cookies owns HTTP cookie encoding only. Authentication, rate limits,
// Origin checks, and CSRF enforcement are implemented in Launcher gateway.
package cookies

import (
	"net/http"
	"strings"

	coregateway "grovehub/libs/core/gateway"
)

const (
	AccessCookie = coregateway.AccessCookie
	RefreshCookie = coregateway.RefreshCookie
	CSRFName = coregateway.CSRFName
	ResetCookie = coregateway.ResetCookie
	SetupCookie = coregateway.SetupCookie
)

const refreshPath = "/identity/public/auth/refresh"

func IssueCookies(writer http.ResponseWriter, access string, refresh string, csrf string) {
	http.SetCookie(writer, &http.Cookie{
		Name: AccessCookie,
		Value: access,
		Path: "/",
		MaxAge: 900,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(writer, &http.Cookie{
		Name: RefreshCookie,
		Value: refresh,
		Path: refreshPath,
		MaxAge: 30 * 24 * 60 * 60,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	if csrf != "" {
		http.SetCookie(writer, &http.Cookie{
			Name: CSRFName,
			Value: csrf,
			Path: "/",
			MaxAge: 30 * 24 * 60 * 60,
			Secure: true,
			SameSite: http.SameSiteStrictMode,
		})
	}
	writer.Header().Set("Cache-Control", "no-store")
}

func ClearSessionCookies(writer http.ResponseWriter) {
	paths := []struct {
		name string
		path string
	}{
		{AccessCookie, "/"},
		{RefreshCookie, refreshPath},
		{CSRFName, "/"},
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
	writer.Header().Set("Cache-Control", "no-store")
}

func ReadCookie(request *http.Request, name string) string {
	cookie, err := request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func SetGrant(writer http.ResponseWriter, name string, path string, value string) {
	http.SetCookie(writer, &http.Cookie{
		Name: name,
		Path: path,
		Value: value,
		MaxAge: 600,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	writer.Header().Set("Cache-Control", "no-store")
}

func ClearGrant(writer http.ResponseWriter, name string, path string) {
	http.SetCookie(writer, &http.Cookie{
		Name: name,
		Path: path,
		MaxAge: -1,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	writer.Header().Set("Cache-Control", "no-store")
}

func Proof(request *http.Request, prefix string) string {
	authorization := request.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, prefix+" ") {
		return ""
	}
	return strings.TrimPrefix(authorization, prefix+" ")
}
