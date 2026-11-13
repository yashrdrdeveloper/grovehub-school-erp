package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sourceClient = `package client
import (
"context"
)
// @Client personClient
type Client interface {
	// @RequestLine POST /protected/persons
	CreatePerson(ctx context.Context, request CreateDTO) (Person, bool, error)
	// @RequestLine GET /protected/persons/{gid}
	GetPersonByGID(parent context.Context, request ByGID) (Person, error)
	// @RequestLine PATCH /protected/persons/{gid}
	UpdatePerson(ctx context.Context, request UpdateDTO) (Person, error)
	// @RequestLine POST /protected/persons/search
	SearchPersons(ctx context.Context, request SearchDTO) (SearchResult, error)
	// @RequestLine PATCH /protected/persons/{gid}/deactivate
	DeactivatePerson(ctx context.Context, request ByGID) (Person, error)
}
`

func fixture(t *testing.T, clientContent, prefix string) (string,string,string) {
	t.Helper()
	root := t.TempDir()
	owner := filepath.Join(root,"services","foundation")
	client := filepath.Join(owner,"client","person.go")
	routes := filepath.Join(owner,".generated","routes","zz_routes.gen.go")
	output := filepath.Join(owner,".generated","client","impl","zz_clients.gen.go")
	for path, content := range map[string]string{
		client: clientContent,
		routes: `package routes
import "net/http"
func Register(mux *http.ServeMux) {
mux.Handle("POST `+prefix+`/protected/persons", nil)
mux.Handle("GET `+prefix+`/protected/persons/{gid}", nil)
mux.Handle("PATCH `+prefix+`/protected/persons/{gid}", nil)
mux.Handle("POST `+prefix+`/protected/persons/search", nil)
mux.Handle("PATCH `+prefix+`/protected/persons/{gid}/deactivate", nil)
}
`,
	} {
		if err := os.MkdirAll(filepath.Dir(path),0755);err!=nil{t.Fatal(err)}
		if err := os.WriteFile(path,[]byte(content),0644);err!=nil{t.Fatal(err)}
	}
	return root, client, output
}

func TestRequestLineGeneratesFullControllerDispatch(t *testing.T) {
	root,_,out := fixture(t,sourceClient,"/foundation")
	if err:=generateAll(root);err!=nil{t.Fatal(err)}
	data,err:=os.ReadFile(out);if err!=nil{t.Fatal(err)}
	generated:=string(data)
	for _,expected:=range []string{
		"httpendpoint.CallStruct[Person](ctx, c.handler, \"POST\", \"/foundation/protected/persons\", request)",
		"httpendpoint.CallStruct[Person](parent, c.handler, \"GET\", \"/foundation/protected/persons/{gid}\", request)",
		"httpendpoint.CallStruct[Person](ctx, c.handler, \"PATCH\", \"/foundation/protected/persons/{gid}\", request)",
		"httpendpoint.CallStruct[SearchResult](ctx, c.handler, \"POST\", \"/foundation/protected/persons/search\", request)",
		"status == 201 && err == nil",
	} {
		if !strings.Contains(generated,expected){t.Errorf("generated client missing %q\n%s",expected,generated)}
	}
	if err:=generateAll(root);err!=nil{t.Fatal(err)}
	again,_:=os.ReadFile(out)
	if string(again)!=generated{t.Fatal("generation must be deterministic")}
}

func TestServicePrefixChangesOnlyInGeneratedRoutes(t *testing.T) {
	root,_,out:=fixture(t,sourceClient,"/school")
	if err:=generateAll(root);err!=nil{t.Fatal(err)}
	data,_:=os.ReadFile(out)
	if !strings.Contains(string(data),"\"/school/protected/persons/{gid}\""){
		t.Fatal("client did not inherit service prefix from generated controller routes")
	}
	if strings.Contains(string(data),"\"/foundation/"){t.Fatal("stale service path")}
}

func TestRequestLineMissingControllerRouteFailsClosed(t *testing.T) {
	root,_,_:=fixture(t,strings.Replace(sourceClient,"GET /protected/persons/{gid}","DELETE /protected/persons/{gid}",1),"/foundation")
	if err:=generateAll(root);err==nil||!strings.Contains(err.Error(),"no matching generated controller route"){
		t.Fatalf("expected route mismatch error; got %v",err)
	}
}

func TestRequestLineRejectsAmbiguousMatches(t *testing.T) {
	root,_,_:=fixture(t,sourceClient,"/foundation")
	routes:=filepath.Join(root,"services","foundation",".generated","routes","zz_routes.gen.go")
	data,_:=os.ReadFile(routes)
	src:=strings.Replace(string(data), "mux.Handle(\"GET /foundation/protected/persons/{gid}\", nil)", "mux.Handle(\"GET /foundation/protected/persons/{gid}\", nil)\nmux.Handle(\"GET /new/protected/persons/{gid}\", nil)",1)
	if err:=os.WriteFile(routes,[]byte(src),0644);err!=nil{t.Fatal(err)}
	if err:=generateAll(root);err==nil||!strings.Contains(err.Error(),"matches multiple routes"){
		t.Fatalf("expected ambiguous matching error; got %v",err)
	}
}

func TestRequestLineRejectsInvalidSignatures(t *testing.T) {
	cases:=[]struct{source,errorFragment string}{
		{strings.Replace(sourceClient,"(Person, bool, error)","(Person, string, error)",1),"second result must be bool"},
		{strings.Replace(sourceClient,"request ByGID","request []ByGID",1),"second @RequestLine parameter must be a named request struct"},
		{strings.Replace(sourceClient,"// @RequestLine GET /protected/persons/{gid}","// @RequestLine GET /protected/persons/{gid}\n\t// @Proxy service.Get",1),"exactly one of"},
		{strings.Replace(sourceClient,"// @RequestLine GET /protected/persons/{gid}","// @RequestLine GET /protected/persons/{gid}\n\t// @RequestLine GET /protected/persons/{gid}",1),"duplicate @RequestLine"},
	}
	for _,item:=range cases{
		root,_,_:=fixture(t,item.source,"/foundation")
		if err:=generateAll(root);err==nil||!strings.Contains(err.Error(),item.errorFragment){
			t.Fatalf("expected %q; got %v",item.errorFragment,err)
		}
	}
}

func TestLegacyProxyMethodsRemainSupported(t *testing.T) {
	client:=`package client
import "context"
// @Client existingClient
type Client interface {
	// @Proxy personService.Get
	Get(ctx context.Context, id string) (string,error)
}
`
	root,_,out:=fixture(t,client,"/foundation")
	if err:=generateAll(root);err!=nil{t.Fatal(err)}
	b,_:=os.ReadFile(out)
	if !strings.Contains(string(b),"return c.personService.Get(ctx, id)") {t.Fatal("legacy @Proxy forwarder was changed")}
	if strings.Contains(string(b),"httpendpoint"){t.Fatal("unused endpoint import")}
}
