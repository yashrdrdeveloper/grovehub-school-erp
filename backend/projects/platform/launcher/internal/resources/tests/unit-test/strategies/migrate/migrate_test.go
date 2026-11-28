package migrate_test

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"gorm.io/gorm"

	"grovehub/launcher/internal/src/strategies/migrate"
	"grovehub/libs/platform/bootstrap"
)

type database struct {
	businessError, auditError bool
	opened, closed            []string
}

func (d *database) Business(context.Context) (*gorm.DB, func(), error) {
	d.opened = append(d.opened, "business")
	if d.businessError {
		return nil, nil, errors.New("business")
	}

	return &gorm.DB{}, func() { d.closed = append(d.closed, "business") }, nil
}

func (d *database) Audit(_ context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
	name := "audit"
	if shared != nil {
		name = "shared-audit"
	}
	d.opened = append(d.opened, name)
	if d.auditError {
		return nil, nil, errors.New("audit")
	}

	return &gorm.DB{}, func() { d.closed = append(d.closed, name) }, nil
}

func (d *database) CurrentName(context.Context, *gorm.DB) (string, error) {
	panic("migration must not use operator name verification")
}

func TestMigration(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		args                      []string
		businessError, auditError bool
		fail                      string
		code                      int
		ran, closed               []string
	}{
		{name: "missing", code: 2},
		{name: "invalid", args: []string{"bad"}, code: 2},
		{name: "business open", args: []string{"foundation"}, businessError: true, code: 1},
		{name: "single", args: []string{"identity"}, ran: []string{"identity"}, closed: []string{"business"}},
		{name: "all", args: []string{"all"}, ran: []string{"foundation", "infrastructure", "integration", "identity", "audit"}, closed: []string{"shared-audit", "business"}},
		{name: "business failure", args: []string{"all"}, fail: "infrastructure", code: 1, ran: []string{"foundation", "infrastructure"}, closed: []string{"business"}},
		{name: "audit only", args: []string{"audit"}, ran: []string{"audit"}, closed: []string{"audit"}},
		{name: "audit open", args: []string{"audit"}, auditError: true, code: 1},
		{name: "audit failure", args: []string{"audit"}, fail: "audit", code: 1, ran: []string{"audit"}, closed: []string{"audit"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := &database{businessError: tc.businessError, auditError: tc.auditError}
			var ran []string
			var operations []migrate.Target
			for _, service := range []string{"foundation", "infrastructure", "integration", "identity", "audit"} {
				target := bootstrap.BusinessDatabase
				if service == "audit" {
					target = bootstrap.AuditDatabase
				}
				operations = append(operations, migrate.Target{
					Name: service,
					Migration: bootstrap.Migration[*gorm.DB]{
						Target: target,
						Run: func(context.Context, *gorm.DB) error {
							ran = append(ran, service)
							if tc.fail == service {
								return errors.New("migration")
							}

							return nil
						},
					},
				})
			}

			code := migrate.New(db, operations).Run(tc.args, io.Discard, io.Discard)
			if code != tc.code || !reflect.DeepEqual(ran, tc.ran) || !reflect.DeepEqual(db.closed, tc.closed) {
				t.Fatalf("code=%d ran=%v closed=%v", code, ran, db.closed)
			}
		})
	}
}
