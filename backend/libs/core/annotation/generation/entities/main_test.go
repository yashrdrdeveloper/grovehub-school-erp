package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateEntityCompatibility(t *testing.T) {
	root:=t.TempDir()
	src:=filepath.Join(root,"entity")
	out:=filepath.Join(root,".generated","entity")
	if err:=os.MkdirAll(src,0755);err!=nil{t.Fatal(err)}
	code:=`package entity
import "time"

// @Entity example.contacts
type Contact struct {
    ID string
    CreatedAt time.Time
}

// @Entity example.persons
type Person struct {
    ID string
    Status string
}
func (p Person) GetID() string {return p.ID}
func (p Person) SetStatus(status string) Person {p.Status=status;return p}
`
	if err:=os.WriteFile(filepath.Join(src,"models.go"),[]byte(code),0644);err!=nil{t.Fatal(err)}
	if err:=generate(src,out);err!=nil{t.Fatal(err)}
	generated,err:=os.ReadFile(filepath.Join(out,"models_gombok.go"));if err!=nil{t.Fatal(err)}
	parts:=[]string{
		"type ContactBuilder struct",
		"func NewContactBuilder() ContactBuilder",
		"func (b ContactBuilder) WithCreatedAt(value time.Time) ContactBuilder",
		"func (b ContactBuilder) Build() Contact",
		"func (value *Contact) GetID() string",
		"func (value *Contact) SetID(next string)",
		"func NewPersonBuilder() PersonBuilder",
		"func (b PersonBuilder) WithStatus(value string) PersonBuilder",
	}
	for _,part:=range parts {
		if !strings.Contains(string(generated),part){t.Errorf("missing generated API %q",part)}
	}
	for _,part:=range []string{"func (value *Person) GetID","func (value *Person) SetStatus"} {
		if strings.Contains(string(generated),part){t.Errorf("collision with hand-written Person accessors: %s",part)}
	}
	if err:=generate(src,out);err!=nil{t.Fatal(err)}
	again,err:=os.ReadFile(filepath.Join(out,"models_gombok.go"));if err!=nil{t.Fatal(err)}
	if string(generated)!=string(again){t.Fatal("entity generation is not deterministic")}
}

func TestRemoveStaleGeneratedWhenAnnotationRemoved(t *testing.T) {
	root:=t.TempDir()
	src:=filepath.Join(root,"entity")
	out:=filepath.Join(root,".generated","entity")
	if err:=os.MkdirAll(src,0755);err!=nil{t.Fatal(err)}
	file:=filepath.Join(src,"entity.go")
	if err:=os.WriteFile(file,[]byte("package entity\n// @Builder\ntype Sample struct {ID string}\n"),0644);err!=nil{t.Fatal(err)}
	if err:=generate(src,out);err!=nil{t.Fatal(err)}
	name:=filepath.Join(out,"entity_gombok.go")
	if _,err:=os.Stat(name);err!=nil{t.Fatal(err)}
	if err:=os.WriteFile(file,[]byte("package entity\ntype Sample struct {ID string}\n"),0644);err!=nil{t.Fatal(err)}
	if err:=generate(src,out);err!=nil{t.Fatal(err)}
	if _,err:=os.Stat(name);!os.IsNotExist(err){t.Fatalf("stale generated output remains: %v",err)}
}

func TestGenerateFoundationContracts(t *testing.T) {
	backend:=filepath.Join("..","..","..","..","..")
	foundation:=filepath.Join(backend,"services","foundation","entity")
	stage:=t.TempDir()
	files:=[]string{"contact.go","person.go","gid.go"}
	for _,file:=range files {
		src,err:=os.ReadFile(filepath.Join(foundation,file));if err!=nil{t.Fatal(err)}
		if err:=os.WriteFile(filepath.Join(stage,file),src,0644);err!=nil{t.Fatal(err)}
	}
	out:=filepath.Join(t.TempDir(),"entity")
	if err:=generate(stage,out);err!=nil{t.Fatal(err)}
	contact,err:=os.ReadFile(filepath.Join(out,"contact_gombok.go"));if err!=nil{t.Fatal(err)}
	person,err:=os.ReadFile(filepath.Join(out,"person_gombok.go"));if err!=nil{t.Fatal(err)}
	if !strings.Contains(string(contact),"GetValue")||!strings.Contains(string(contact),"SetValue"){t.Fatal("Contact accessors missing")}
	if strings.Contains(string(person),"func (value *Person) SetStatus"){t.Fatal("Person transition method was overridden")}
}


func TestEntityTableOnlyDoesNotGenerateAccessors(t *testing.T) {
    source:=t.TempDir()
    output:=filepath.Join(t.TempDir(),"entity")
    code:="package entity\n// @Entity identity.password_credentials mode=table\ntype PasswordCredential struct {PasswordHash string}\n"
    if err:=os.WriteFile(filepath.Join(source,"password.go"),[]byte(code),0644);err!=nil{t.Fatal(err)}
    if err:=generate(source,output);err!=nil{t.Fatal(err)}
    if _,err:=os.Stat(filepath.Join(output,"password_gombok.go"));!os.IsNotExist(err){t.Fatalf("table-only entity generated accessors: %v",err)}
}
