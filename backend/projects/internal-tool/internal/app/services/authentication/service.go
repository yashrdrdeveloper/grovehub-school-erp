package authentication

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	corehttp "grovehub/libs/core/http"
	"log/slog"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
	operators "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/operator"
	sessions "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/session"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalid = errors.New("invalid credentials or session")
	ErrLimited = errors.New("authentication rate limited")
)

const SessionTTL = 8 * time.Hour

type Service struct {
	operators *operators.Repository
	sessions  *sessions.Repository
	hashing   chan struct{}
}

func New(db *sql.DB) *Service {
	return &Service{operators.New(db), sessions.New(db), make(chan struct{}, 2)}
}

func normalizeEmail(email string) (string, bool) {
	email = strings.ToLower(strings.TrimSpace(email))
	parsed, err := mail.ParseAddress(email)
	return email, err == nil && parsed.Address == email && len(email) <= 254
}
func hashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 32)
}
func digest(raw string) []byte { sum := sha256.Sum256([]byte(raw)); return sum[:] }
func secret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func validSecret(raw string) bool {
	if len(raw) != 43 {
		return false
	}
	b, err := base64.RawURLEncoding.DecodeString(raw)
	return err == nil && len(b) == 32 && base64.RawURLEncoding.EncodeToString(b) == raw
}

func (s *Service) Bootstrap(ctx context.Context, email, password string) error {
	email, ok := normalizeEmail(email)
	if !ok || len(password) < 12 || len(password) > 256 {
		return ErrInvalid
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	return s.operators.Bootstrap(ctx, entity.Credential{Operator: entity.Operator{ID: uuid.NewString(), Email: email, Persona: entity.L2}, Salt: salt, Hash: hashPassword(password, salt)})
}

type LoginResult struct {
	Operator    entity.Operator
	Token, CSRF string
}

func (s *Service) Login(ctx context.Context, email, password string) (result LoginResult, loginErr error) {
	ctx = corehttp.Ensure(ctx)
	defer func() {
		outcome := ""
		if errors.Is(loginErr, ErrInvalid) {
			outcome = "REJECTED"
		} else if errors.Is(loginErr, ErrLimited) {
			outcome = "LIMITED"
		}
		if outcome == "" {
			return
		}
		record, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()
		if err := s.sessions.CountAuthenticationOutcome(record, outcome); err != nil {
			slog.Warn("authentication counter unavailable")
		}
	}()
	email, ok := normalizeEmail(email)
	if !ok || len(password) == 0 || len(password) > 256 {
		return LoginResult{}, ErrInvalid
	}
	select {
	case s.hashing <- struct{}{}:
		defer func() { <-s.hashing }()
	default:
		return LoginResult{}, ErrLimited
	}
	allowed, err := s.sessions.ClaimLoginAttempt(ctx, digest(email))
	if err != nil {
		return LoginResult{}, err
	}
	if !allowed {
		return LoginResult{}, ErrLimited
	}
	credential, err := s.operators.Credential(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return LoginResult{}, err
	}
	// Unknown and disabled accounts still incur the same bounded Argon2 work.
	salt := make([]byte, 16)
	expected := make([]byte, 32)
	if err == nil {
		salt = credential.Salt
		expected = credential.Hash
	}
	match := subtle.ConstantTimeCompare(hashPassword(password, salt), expected) == 1
	if err != nil || !credential.Active || !match {
		return LoginResult{}, ErrInvalid
	}
	token, err := secret()
	if err != nil {
		return LoginResult{}, err
	}
	csrf, err := secret()
	if err != nil {
		return LoginResult{}, err
	}
	if err = s.sessions.Create(ctx, credential.Operator.ID, digest(token), digest(csrf), time.Now().Add(SessionTTL)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResult{}, ErrInvalid
		}
		return LoginResult{}, err
	}
	return LoginResult{credential.Operator, token, csrf}, nil
}
func (s *Service) Verify(ctx context.Context, token string) (entity.Operator, error) {
	if !validSecret(token) {
		return entity.Operator{}, ErrInvalid
	}
	op, _, err := s.sessions.Verify(ctx, digest(token))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Operator{}, ErrInvalid
	}
	return op, err
}
func (s *Service) VerifyMutation(ctx context.Context, token, csrf string) (entity.Operator, error) {
	if !validSecret(token) || !validSecret(csrf) {
		return entity.Operator{}, ErrInvalid
	}
	op, expected, err := s.sessions.Verify(ctx, digest(token))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Operator{}, ErrInvalid
	}
	if err != nil {
		return entity.Operator{}, err
	}
	if subtle.ConstantTimeCompare(digest(csrf), expected) != 1 {
		return entity.Operator{}, ErrInvalid
	}
	return op, nil
}
func (s *Service) Logout(ctx context.Context, token, csrf string) error {
	if _, err := s.VerifyMutation(ctx, token, csrf); err != nil {
		return err
	}
	return s.sessions.Revoke(ctx, digest(token))
}
