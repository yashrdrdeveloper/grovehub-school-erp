package postgres

import (
	"context"
	"os"
	"testing"
)

func TestProjectConfig(t *testing.T) {
	for _, dsn := range []string{"", "host=localhost dbname=school", "postgres://localhost/school", "postgres://school@localhost", "postgres://school@localhost/school?dbname=platform", "postgres://school@localhost/school?user=platform", "postgres://school@localhost/school?service=platform"} {
		if _, err := projectConfig(dsn); err == nil { t.Errorf("accepted invalid URI %q", dsn) }
	}
	config, err := projectConfig("postgres://school:password@127.0.0.1:5432/school?sslmode=disable")
	if err != nil || config.Database != "school" || config.User != "school" { t.Fatal("explicit project configuration rejected") }
}

func TestProjectDatabaseIsolation(t *testing.T) {
	school := os.Getenv("PROJECT_SCHOOL_TEST_DSN")
	internal := os.Getenv("PROJECT_INTERNAL_TEST_DSN")
	denied := os.Getenv("PROJECT_CROSS_TEST_DSN")
	if school == "" && internal == "" && denied == "" { t.Skip("requires disposable project databases") }
	if school == "" || internal == "" || denied == "" { t.Fatal("incomplete integration test configuration") }
	for _, dsn := range []string{school, internal} {
		db, err := OpenProjectPool(context.Background(), dsn)
		if err != nil { t.Fatal(err) }
		if err := db.PingContext(context.Background()); err != nil { t.Fatal("project pool unavailable") }
		_ = db.Close()
	}
	if db, err := OpenProjectPool(context.Background(), denied); err == nil { _ = db.Close(); t.Fatal("cross-project role unexpectedly connected") }
}
