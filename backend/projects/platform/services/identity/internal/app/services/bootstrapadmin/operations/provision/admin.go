// Command bootstrap-admin provisions the first Identity administrator from an
// operator-controlled terminal. It is never registered as an HTTP route.
package bootstrapadmin

import (
    "bufio"
    "context"
    "errors"
    "flag"
    "fmt"
    "net/mail"
    "io"
    "os"
    "strings"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"

    coreid "grovehub/libs/core/id"
    coredatabase "grovehub/libs/core/database"
    passwords "grovehub/services/identity/internal/app/services/credential/helpers"
)

// Run is a privileged operator entrypoint, never an HTTP handler.
func Run() error {
	return RunArgs(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}

func RunArgs(args []string, input io.Reader, output, diagnostics io.Writer) error {
    flags := flag.NewFlagSet("identity-bootstrap-admin", flag.ContinueOnError)
    flags.SetOutput(diagnostics)
    var email, first, last, tenant, personGID string
    var passwordStdin, confirm bool
    flags.StringVar(&email, "email", "", "Initial administrator email")
    flags.StringVar(&first, "first-name", "", "First name for a new Foundation Person")
    flags.StringVar(&last, "last-name", "", "Last name for a new Foundation Person")
    flags.StringVar(&tenant, "tenant-id", "", "Non-empty tenant scope")
    flags.StringVar(&personGID, "person-gid", "", "Optional existing ACTIVE Foundation Person UUID")
    flags.BoolVar(&passwordStdin, "password-stdin", false, "Read the password from stdin (never from flags)")
    flags.BoolVar(&confirm, "confirm-bootstrap", false, "Explicitly authorize one-time provisioning")
    if err := flags.Parse(args); err != nil { return err }

    if !confirm || !passwordStdin || len(flags.Args()) != 0 {
        return errors.New("required: --confirm-bootstrap --password-stdin --email --tenant-id (and --first-name/--last-name for a new person)")
    }
    email = strings.ToLower(strings.TrimSpace(email))
    address, err := mail.ParseAddress(email)
    if err != nil || address.Address != email || address.Name != "" || len(email) > 254 {
        return errors.New("invalid email")
    }
    tenant = strings.TrimSpace(tenant)
    if tenant == "" || len(tenant) > 128 {
        return errors.New("tenant-id must be 1-128 characters")
    }
    if personGID != "" {
        parsed, err := uuid.Parse(personGID)
        if err != nil { return errors.New("invalid person-gid UUID") }
        personGID = parsed.String()
    } else {
        first, last = strings.TrimSpace(first), strings.TrimSpace(last)
        if first == "" || last == "" || len(first) > 100 || len(last) > 100 {
            return errors.New("first-name and last-name must be 1-100 characters")
        }
    }
    if os.Getenv("DATABASE_URL") == "" {
        return errors.New("DATABASE_URL is missing")
    }

    reader := bufio.NewReaderSize(input, 512)
    password, err := reader.ReadString('\n')
    if err != nil && !errors.Is(err, io.EOF) {
        return errors.New("could not read password from stdin")
    }
    password = strings.TrimSuffix(strings.TrimSuffix(password, "\n"), "\r")
    if len(password) < 12 || len(password) > 256 {
        return errors.New("password must be between 12 and 256 bytes")
    }
    hash, err := passwords.HashPassword(password)
    password = ""
    if err != nil { return fmt.Errorf("password hashing failed: %w", err) }

    db, err := coredatabase.NewDB(os.Getenv("DATABASE_URL"))
    if err != nil { return fmt.Errorf("database connection: %w", err) }
    sqlDB, err := db.DB()
    if err != nil { return err }
    defer sqlDB.Close()
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    accountID, err := coreid.NewUUIDv7()
    if err != nil { return err }
    err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Serializes concurrent bootstrap commands on the same PostgreSQL DB.
        if err := tx.Exec("SELECT pg_advisory_xact_lock(743029, 119033)").Error; err != nil { return err }
        var count int64
        if err := tx.Raw("SELECT COUNT(*) FROM identity.accounts").Scan(&count).Error; err != nil { return err }
        if count != 0 { return errors.New("bootstrap refused: Identity already has accounts; use normal administrator provisioning") }

        if personGID == "" {
            // Reuse an existing Foundation Person with the same normalized email,
            // otherwise create Person and primary email in this bootstrap transaction.
            var ids []string
            if err := tx.Raw("SELECT person_gid::text FROM foundation.person_contacts WHERE type = 'EMAIL' AND normalized_value = ?", email).Scan(&ids).Error; err != nil { return err }
            if len(ids) > 1 { return errors.New("ambiguous Foundation email") }
            if len(ids) == 1 {
                personGID = ids[0]
            } else {
                personGID, err = coreid.NewUUIDv7()
                if err != nil { return err }
                now := time.Now().UTC()
                search := strings.ToLower(strings.Join([]string{first, last, email}, " "))
                if err := tx.Exec("INSERT INTO foundation.persons (gid, first_name, last_name, status, version, search_details, created_at, updated_at) VALUES (?, ?, ?, 'ACTIVE', 1, ?, ?, ?)", personGID, first, last, search, now, now).Error; err != nil { return err }
                if err := tx.Exec("INSERT INTO foundation.person_contacts (person_gid, type, value, normalized_value, is_primary, is_verified, created_at, updated_at) VALUES (?, 'EMAIL', ?, ?, TRUE, FALSE, ?, ?)", personGID, email, email, now, now).Error; err != nil { return err }
            }
        }
        var valid int64
        if err := tx.Raw("SELECT COUNT(*) FROM foundation.persons WHERE gid = ? AND status = 'ACTIVE'", personGID).Scan(&valid).Error; err != nil { return err }
        if valid != 1 { return errors.New("Foundation Person must exist and be ACTIVE") }
        var matching int64
        if err := tx.Raw("SELECT COUNT(*) FROM foundation.person_contacts WHERE person_gid = ? AND type = 'EMAIL' AND normalized_value = ?", personGID, email).Scan(&matching).Error; err != nil { return err }
        if matching != 1 { return errors.New("email must belong to the selected Foundation Person") }

        now := time.Now().UTC()
        if err := tx.Exec("INSERT INTO identity.accounts (account_id, person_gid, tenant_id, email, status, version, search_details, created_at, updated_at) VALUES (?, ?, ?, ?, 'ACTIVE', 1, ?, ?, ?)", accountID, personGID, tenant, email, email, now, now).Error; err != nil { return err }
        if err := tx.Exec("INSERT INTO identity.password_credentials (account_id, password_hash, password_changed_at, version, created_at, updated_at) VALUES (?, ?,  ?, 1, ?, ?)", accountID, hash, now, now, now).Error; err != nil { return err }
        return nil
    })
    if err != nil { return fmt.Errorf("transaction rolled back: %w", err) }
    fmt.Fprintf(output, "First account created successfully.\nAccount ID: %s\nTenant ID: %s\n", accountID, tenant)
    fmt.Fprintln(output, "To grant admin privileges set IDENTITY_ADMIN_ACCOUNT_IDS to the Account ID and RESTART the backend.")
    fmt.Fprintln(output, "No password or session secret was printed.")
    return nil
}
