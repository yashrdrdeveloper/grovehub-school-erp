package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateStableEventContract(t *testing.T) {
	backend:=t.TempDir()
	dir:=filepath.Join(backend,"services","foundation","events","persons")
	if err:=os.MkdirAll(dir,0755);err!=nil{t.Fatal(err)}
	source:=`package persons
// @Event topic=foundation_person_changed type=foundation.person.changed.v1 source=foundation
type Changed struct {ID string}
`
	if err:=os.WriteFile(filepath.Join(dir,"changed.go"),[]byte(source),0644);err!=nil{t.Fatal(err)}
	if err:=generateAll(backend);err!=nil{t.Fatal(err)}
	output:=filepath.Join(backend,"services","foundation",".generated","events","persons","zz_events.gen.go")
	data,err:=os.ReadFile(output);if err!=nil{t.Fatal(err)}
	for _,part:=range []string{"ChangedTopic",`"foundation_person_changed"`,"ChangedType",`"foundation.person.changed.v1"`,"ChangedVersion = 1",`Source`, `"foundation"`} {
		if !strings.Contains(string(data),part){t.Errorf("event contract missing %q: %s",part,data)}
	}
	if err:=generateAll(backend);err!=nil{t.Fatal(err)}
	again,err:=os.ReadFile(output);if err!=nil{t.Fatal(err)}
	if string(data)!=string(again){t.Fatal("nondeterministic event output")}
}

func TestRejectInvalidEventType(t *testing.T) {
	dir:=t.TempDir()
	source:=`package events
// @Event topic=domain_entity type=domain.entity.invalid
type Changed struct{}
`
	if err:=os.WriteFile(filepath.Join(dir,"changed.go"),[]byte(source),0644);err!=nil{t.Fatal(err)}
	_,_,err:=discover(dir)
	if err==nil||!strings.Contains(err.Error(),"requires valid topic"){t.Fatalf("invalid event version accepted: %v",err)}
}

func TestCurrentEventContractsPreserveNames(t *testing.T) {
	backend:=filepath.Join("..","..","..","..","..")
	tests:=[]struct{service,path,topic,eventType string}{
		{"foundation","persons","foundation_person_changed","foundation.person.changed.v1"},
		{"identity","accounts","identity_account_created","identity.account.created.v1"},
		{"integration","connections/mailbox","integration_mailbox_changed","integration.mailbox.changed.v1"},
	}
	for _,tc:=range tests {
		dir:=filepath.Join(backend,"services",tc.service,"events",filepath.FromSlash(tc.path))
		defs,_,err:=discover(dir);if err!=nil{t.Fatalf("%s: %v",tc.service,err)}
		if len(defs)!=1||defs[0].topic!=tc.topic||defs[0].eventType!=tc.eventType||defs[0].version!=1 {
			t.Errorf("%s: event contract drift: %+v",tc.service,defs)
		}
	}
}
