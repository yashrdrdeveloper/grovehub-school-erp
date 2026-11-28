package readiness_test

import (
	corehttp "grovehub/libs/core/http"

	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestHandler(t *testing.T) {
	var drain atomic.Bool
	failed := false
	check := func(context.Context) error {
		if failed { return errors.New("private database detail") }
		return nil
	}
	h := corehttp.ReadinessHandler(&drain, check)
	for _, step := range []struct{ fail, draining bool; want int }{
		{false,false,204},{true,false,503},{false,true,503},
	} {
		failed = step.fail
		drain.Store(step.draining)
		w := httptest.NewRecorder()
		h(w,httptest.NewRequest("GET","/health/ready",nil))
		if w.Code != step.want || w.Body.Len()!=0 || w.Header().Get("Cache-Control")!="no-store" { t.Fatalf("unexpected readiness response: %d %s",w.Code,w.Body) }
	}
}

func TestHTTP(t *testing.T) {
	var status atomic.Int32
	status.Store(200)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.URL.Path!="/ready" || r.Header.Get("Authorization")!="" || r.Header.Get("Cookie")!="" { t.Error("unexpected health request") }
		if status.Load()==302 { w.Header().Set("Location","/ready") }
		w.WriteHeader(int(status.Load()))
	}))
	defer server.Close()
	check,err:=corehttp.ReadinessHTTP(server.URL,"/ready")
	if err!=nil {t.Fatal(err)}
	if err=check(context.Background());err!=nil {t.Fatal(err)}
	for _,code:=range []int32{302,401,503}{
		status.Store(code)
		if err=check(context.Background());err==nil {t.Fatalf("accepted %d",code)}
	}
	for _,origin:=range []string{"http://example.com","https://example.com/path","https://user@example.com","https://example.com?","https://example.com#fragment"} {
		if _,err=corehttp.ReadinessHTTP(origin,"/ready");err==nil {t.Fatalf("accepted %s",origin)}
	}
}
