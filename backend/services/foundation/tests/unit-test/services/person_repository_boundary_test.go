package services_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	coresearch "grovehub/libs/core/search"
	annotationvalidation "grovehub/libs/core/annotation/validation"
	annotationdto "grovehub/libs/core/annotation/dto"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
	repositoryerrors "grovehub/services/foundation/repository/constants/errors"
	personservice "grovehub/services/foundation/services/impl"
)

type transactionMarker struct{}

type fakeUnitOfWork struct {
	runs int
}

func (u *fakeUnitOfWork) Run(ctx context.Context, action func(context.Context) error) error {
	u.runs++
	return action(context.WithValue(ctx, transactionMarker{}, true))
}

type fakePersonRepository struct {
	person  entity.Person
	log     *[]string
	queries int
}

func (r *fakePersonRepository) Create(ctx context.Context, person entity.Person) error {
	if ctx.Value(transactionMarker{}) != true {
		return errors.New("person create was not inside unit of work")
	}
	*r.log = append(*r.log, "person.create")
	r.person = person
	return nil
}

func (r *fakePersonRepository) FindByGID(_ context.Context, gid string) (entity.Person, error) {
	if r.person.GID != gid {
		return entity.Person{}, repositoryerrors.ErrPersonNotFound
	}
	return r.person, nil
}

func (r *fakePersonRepository) FindByIdempotencyKey(_ context.Context, key string) (entity.Person, error) {
	if r.person.IdempotencyKey != nil && *r.person.IdempotencyKey == key {
		return r.person, nil
	}
	return entity.Person{}, repositoryerrors.ErrPersonNotFound
}

func (r *fakePersonRepository) Search(_ context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error) {
	r.queries++
	return coresearch.QueryResult[entity.Person]{Items: []entity.Person{r.person}}, nil
}

func (r *fakePersonRepository) Update(ctx context.Context, person entity.Person) error {
	if ctx.Value(transactionMarker{}) != true {
		return errors.New("person update was not inside unit of work")
	}
	*r.log = append(*r.log, "person.update")
	r.person = person
	return nil
}

type fakeContactRepository struct {
	byPerson    map[string][]entity.Contact
	log         *[]string
	replacements int
	batchReads   int
}

func (r *fakeContactRepository) FindByPersonGID(_ context.Context, gid string) ([]entity.Contact, error) {
	return append([]entity.Contact(nil), r.byPerson[gid]...), nil
}

func (r *fakeContactRepository) FindByPersonGIDs(_ context.Context, gids []string) (map[string][]entity.Contact, error) {
	r.batchReads++
	result := make(map[string][]entity.Contact, len(gids))
	for _, gid := range gids {
		result[gid] = append([]entity.Contact(nil), r.byPerson[gid]...)
	}
	return result, nil
}

func (r *fakeContactRepository) FindPersonGIDsByContacts(_ context.Context, contacts []entity.Contact) ([]string, error) {
	var result []string
	for gid, stored := range r.byPerson {
		for _, incoming := range contacts {
			for _, known := range stored {
				if known.Type == incoming.Type && known.NormalizedValue == incoming.NormalizedValue {
					result = append(result, gid)
					break
				}
			}
		}
	}
	return result, nil
}

func (r *fakeContactRepository) ReplaceForPerson(ctx context.Context, gid string, contacts []entity.Contact) error {
	if ctx.Value(transactionMarker{}) != true {
		return errors.New("contact replacement was not inside unit of work")
	}
	for _, contact := range contacts {
		if contact.PersonGID != gid {
			return errors.New("contact foreign key does not match person owner")
		}
	}
	*r.log = append(*r.log, "contact.replace")
	r.replacements++
	r.byPerson[gid] = append([]entity.Contact(nil), contacts...)
	return nil
}

type fakeValidator struct{}

func (fakeValidator) CreateRules() annotationvalidation.Rules[foundationrequest.CreatePerson] {
	return annotationvalidation.Rules[foundationrequest.CreatePerson]{}
}
func (fakeValidator) UpdateRules() annotationvalidation.Rules[foundationrequest.UpdatePerson] {
	return annotationvalidation.Rules[foundationrequest.UpdatePerson]{}
}

func (fakeValidator) ValidateCreate(request foundationrequest.CreatePerson) (foundationrequest.CreatePerson, error) {
	return request, nil
}
func (fakeValidator) ValidateUpdate(request foundationrequest.UpdatePerson) (foundationrequest.UpdatePerson, error) {
	return request, nil
}
func (fakeValidator) ValidateGID(string) error { return nil }
func (fakeValidator) ValidateSearch(coresearch.Request) (coresearch.Query[string], coresearch.CursorPayload, error) {
	return coresearch.Query[string]{Sort: coresearch.QuerySort{Field: "gid", Direction: coresearch.DirectionAscending}, Limit: 25}, coresearch.CursorPayload{}, nil
}

type fakeMapper struct{}

func (fakeMapper) ToPersonGID(request foundationrequest.PersonByGID) string { return request.GID }
func (fakeMapper) ToPersonErrorMetadata(gid string) foundationresponse.PersonErrorMetadata {
	return foundationresponse.PersonErrorMetadata{GID: gid}
}
func (fakeMapper) ToPersonResponse(person entity.Person, contacts []entity.Contact) foundationresponse.Person {
	response := foundationresponse.Person{GID: person.GID, FirstName: person.FirstName, Version: person.Version}
	for _, contact := range contacts {
		value := foundationresponse.ContactValue{Value: contact.Value, IsPrimary: contact.IsPrimary, IsVerified: contact.IsVerified}
		if contact.Type == entity.ContactTypeEmail {
			response.Contacts.Emails = append(response.Contacts.Emails, value)
		} else if contact.Type == entity.ContactTypePhone {
			response.Contacts.Phones = append(response.Contacts.Phones, value)
		}
	}
	return response
}

func setup(t *testing.T) (*fakePersonRepository, *fakeContactRepository, *fakeUnitOfWork, interface {
	Create(context.Context, foundationrequest.CreatePerson) (foundationresponse.Person, bool, error)
	Update(context.Context, foundationrequest.PersonByGID, foundationrequest.UpdatePerson) (foundationresponse.Person, error)
	Search(context.Context, coresearch.Request) (coresearch.Result[foundationresponse.Person], error)
}) {
	t.Helper()
	var writes []string
	people := &fakePersonRepository{log: &writes}
	contacts := &fakeContactRepository{byPerson: make(map[string][]entity.Contact), log: &writes}
	work := &fakeUnitOfWork{}
	service := personservice.NewPersonService(people, contacts, work, fakeMapper{}, fakeValidator{})
	return people, contacts, work, service
}

func TestCreatePersistsPersonAndContactsInsideOneTransaction(t *testing.T) {
	people, contacts, work, service := setup(t)
	result, created, err := service.Create(context.Background(), foundationrequest.CreatePerson{
		Person: foundationrequest.Person{FirstName: "Anu", LastName: "Roy"},
		Contacts: foundationrequest.Contacts{
			Emails: []foundationrequest.ContactValue{{Value: "anu@example.com", IsPrimary: true}},
		},
	})
	if err != nil || !created { t.Fatalf("create failed: created=%v, err=%v", created, err) }
	if work.runs != 1 || contacts.replacements != 1 { t.Fatalf("expected one unit of work and one contact replacement; got %d / %d", work.runs, contacts.replacements) }
	if result.GID == "" || people.person.GID != result.GID { t.Fatalf("person GID mismatch: %#v", people.person) }
	if !strings.Contains(people.person.SearchDetails, "anu@example.com") {
		t.Fatalf("Person search projection omits contacts: %q", people.person.SearchDetails)
	}
	if len(result.Contacts.Emails) != 1 || len(contacts.byPerson[result.GID]) != 1 {
		t.Fatalf("contact was not saved and returned: %#v", result.Contacts)
	}
	if got := *people.log; !reflect.DeepEqual(got, []string{"person.create", "contact.replace"}) {
		t.Fatalf("unexpected write order: %v", got)
	}
}

func TestUpdateNameOnlyDoesNotRewriteContacts(t *testing.T) {
	people, contacts, work, service := setup(t)
	people.person = entity.Person{GID: "person-id", FirstName: "Old", LastName: "Name", Version: 1, UpdatedAt: time.Now().UTC()}
	contacts.byPerson[people.person.GID] = []entity.Contact{{PersonGID: people.person.GID, Type: entity.ContactTypeEmail, Value: "old@example.com", NormalizedValue: "old@example.com"}}
	name := "New"
	result, err := service.Update(context.Background(), foundationrequest.PersonByGID{GID: "person-id"}, foundationrequest.UpdatePerson{FirstName: annotationdto.Patch[string]{Set: true, Value: &name}})
	if err != nil { t.Fatal(err) }
	if result.FirstName != "New" || people.person.Version != 2 { t.Fatalf("unexpected person update: %#v", people.person) }
	if contacts.replacements != 0 { t.Fatal("name-only update must not rewrite unchanged contacts") }
	if work.runs != 1 { t.Fatalf("expected a unit of work, got %d", work.runs) }
	if !strings.Contains(people.person.SearchDetails, "old@example.com") {
		t.Fatalf("Person projection lost existing contact: %q", people.person.SearchDetails)
	}
}

func TestUpdateContactsUsesSamePersonUnitOfWork(t *testing.T) {
	people, contacts, work, service := setup(t)
	people.person = entity.Person{GID: "person-id", FirstName: "A", LastName: "B", Version: 2, UpdatedAt: time.Now().UTC()}
	newContacts := foundationrequest.Contacts{Phones: []foundationrequest.ContactValue{{Value: "+919876543210", IsPrimary: true}}}
	result, err := service.Update(context.Background(), foundationrequest.PersonByGID{GID: "person-id"}, foundationrequest.UpdatePerson{Contacts: annotationdto.Patch[foundationrequest.Contacts]{Set: true, Value: &newContacts}})
	if err != nil { t.Fatal(err) }
	if work.runs != 1 || contacts.replacements != 1 || people.person.Version != 3 {
		t.Fatalf("Person and Contact must write in one unit of work: runs=%d replacements=%d version=%d", work.runs, contacts.replacements, people.person.Version)
	}
	if !reflect.DeepEqual(*people.log, []string{"person.update", "contact.replace"}) {
		t.Fatalf("unexpected write order: %v", *people.log)
	}
	if len(result.Contacts.Phones) != 1 || !strings.Contains(people.person.SearchDetails, "+919876543210") {
		t.Fatalf("contact replacement or search projection missing: %#v, %q", result.Contacts, people.person.SearchDetails)
	}
}

func TestSearchLoadsContactRowsFromContactRepository(t *testing.T) {
	people, contacts, work, service := setup(t)
	people.person = entity.Person{GID: "person-id", FirstName: "A", LastName: "B"}
	contacts.byPerson["person-id"] = []entity.Contact{{PersonGID: "person-id", Type: entity.ContactTypeEmail, Value: "a@example.com"}}
	result, err := service.Search(context.Background(), coresearch.Request{})
	if err != nil { t.Fatal(err) }
	if len(result.Items) != 1 || len(result.Items[0].Contacts.Emails) != 1 {
		t.Fatalf("search must combine Person rows with separately loaded Contact rows: %#v", result)
	}
	if people.queries != 1 || contacts.batchReads != 1 || work.runs != 0 {
		t.Fatalf("search unexpectedly mutated data: queries=%d batches=%d transactions=%d", people.queries, contacts.batchReads, work.runs)
	}
}
