package helpers

import (
    "strings"
    "testing"
    "time"

    "grovehub/services/foundation/internal/app/entity"
)

func TestPersonCursorValue(t *testing.T) {
    now := time.Date(2026, 11, 14, 12, 34, 56, 123000000, time.FixedZone("TZ", 5*60*60+30*60))
    person := entity.Person{GID: "p-1", FirstName: "Anu", LastName: "Roy", CreatedAt: now, UpdatedAt: now}
    cases := map[string]string{
        "gid": "p-1",
        "firstName": "anu",
        "lastName": "roy",
        "createdAt": now.UTC().Format(time.RFC3339Nano),
        "updatedAt": now.UTC().Format(time.RFC3339Nano),
        "unsupported": "",
    }
    for field, expected := range cases {
        if got := personCursorValue(person, field); got != expected {
            t.Errorf("%s: got %q want %q", field, got, expected)
        }
    }
    if strings.Contains(personCursorValue(person, "firstName"), "Anu") { t.Fatal("name cursor must be normalized") }
}
