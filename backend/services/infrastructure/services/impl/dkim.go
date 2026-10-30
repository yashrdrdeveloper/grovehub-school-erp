package impl

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"grovehub/services/infrastructure/entity"
	infrastructurerepository "grovehub/services/infrastructure/repository"
	infrastructureservices "grovehub/services/infrastructure/services"
)

const defaultDKIMBits = 2048

type DKIMService struct {
	repository infrastructurerepository.DKIMRepository
	stateDir   string
}

func NewDKIMService(repository infrastructurerepository.DKIMRepository) infrastructureservices.DKIMService {
	stateDir := firstNonEmpty(
		os.Getenv("MAIL_INFRASTRUCTURE_DKIM_STATE_DIR"),
		os.Getenv("MAIL_PROVIDER_DKIM_STATE_DIR"),
		".state/mailProvider/dkim",
	)
	return &DKIMService{repository: repository, stateDir: stateDir}
}

func (service *DKIMService) Ensure(
	ctx context.Context,
	domain string,
	selector string,
) (infrastructureservices.DKIMIdentity, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	selector = strings.TrimSpace(selector)
	if domain == "" || selector == "" {
		return infrastructureservices.DKIMIdentity{}, infrastructureservices.ErrInvalidDKIMIdentity
	}

	persisted, err := service.repository.FindByDomain(ctx, domain)
	if err == nil {
		return service.usePersisted(persisted, selector)
	}
	if !errors.Is(err, infrastructurerepository.ErrDKIMKeyNotFound) {
		return infrastructureservices.DKIMIdentity{}, err
	}

	key, found, err := service.readExistingFilesystemKey(domain, selector)
	if err != nil {
		return infrastructureservices.DKIMIdentity{}, err
	}
	if !found {
		key, err = generateDKIMKey(domain, selector)
		if err != nil {
			return infrastructureservices.DKIMIdentity{}, err
		}
	}

	if err := service.repository.Create(ctx, key); err != nil {
		if !errors.Is(err, infrastructurerepository.ErrDKIMKeyAlreadyExists) {
			return infrastructureservices.DKIMIdentity{}, err
		}
		persisted, readErr := service.repository.FindByDomain(ctx, domain)
		if readErr != nil {
			return infrastructureservices.DKIMIdentity{}, readErr
		}
		return service.usePersisted(persisted, selector)
	}

	if err := service.materialize(key); err != nil {
		return infrastructureservices.DKIMIdentity{}, err
	}
	return identityFromKey(key), nil
}

func (service *DKIMService) usePersisted(
	key entity.DKIMKey,
	selector string,
) (infrastructureservices.DKIMIdentity, error) {
	if key.Selector != selector {
		return infrastructureservices.DKIMIdentity{}, fmt.Errorf(
			"%w: domain %s is permanently bound to selector %s",
			infrastructureservices.ErrDKIMSelectorConflict,
			key.Domain,
			key.Selector,
		)
	}
	if err := service.materialize(key); err != nil {
		return infrastructureservices.DKIMIdentity{}, err
	}
	return identityFromKey(key), nil
}

func (service *DKIMService) readExistingFilesystemKey(
	domain string,
	selector string,
) (entity.DKIMKey, bool, error) {
	keyDir := filepath.Join(service.stateDir, domain)
	privatePEM, err := os.ReadFile(filepath.Join(keyDir, "private"))
	if errors.Is(err, os.ErrNotExist) {
		return entity.DKIMKey{}, false, nil
	}
	if err != nil {
		return entity.DKIMKey{}, false, err
	}
	if len(strings.TrimSpace(string(privatePEM))) == 0 {
		return entity.DKIMKey{}, false, nil
	}

	selectorBytes, selectorErr := os.ReadFile(filepath.Join(keyDir, "selector"))
	if selectorErr == nil {
		existingSelector := strings.TrimSpace(string(selectorBytes))
		if existingSelector != "" && existingSelector != selector {
			return entity.DKIMKey{}, false, fmt.Errorf(
				"%w: filesystem selector %s does not match requested selector %s",
				infrastructureservices.ErrDKIMSelectorConflict,
				existingSelector,
				selector,
			)
		}
	} else if !errors.Is(selectorErr, os.ErrNotExist) {
		return entity.DKIMKey{}, false, selectorErr
	}

	privateKey, err := parseRSAPrivateKey(privatePEM)
	if err != nil {
		return entity.DKIMKey{}, false, fmt.Errorf("parse existing DKIM private key: %w", err)
	}
	publicPEM, err := marshalPublicKey(&privateKey.PublicKey)
	if err != nil {
		return entity.DKIMKey{}, false, err
	}
	return entity.DKIMKey{
		Domain:        domain,
		Selector:      selector,
		PrivateKeyPEM: string(privatePEM),
		PublicKeyPEM:  string(publicPEM),
		CreatedAt:     time.Now().UTC(),
	}, true, nil
}

func (service *DKIMService) materialize(key entity.DKIMKey) error {
	keyDir := filepath.Join(service.stateDir, key.Domain)
	if err := os.MkdirAll(keyDir, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(keyDir, "private"), []byte(key.PrivateKeyPEM), 0o600); err != nil {
		return err
	}
	if err := os.Chmod(filepath.Join(keyDir, "private"), 0o600); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(keyDir, "public"), []byte(key.PublicKeyPEM), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(keyDir, "selector"), []byte(key.Selector+"\n"), 0o644); err != nil {
		return err
	}
	return nil
}

func generateDKIMKey(domain string, selector string) (entity.DKIMKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, defaultDKIMBits)
	if err != nil {
		return entity.DKIMKey{}, err
	}
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	publicPEM, err := marshalPublicKey(&privateKey.PublicKey)
	if err != nil {
		return entity.DKIMKey{}, err
	}
	return entity.DKIMKey{
		Domain:        domain,
		Selector:      selector,
		PrivateKeyPEM: string(privatePEM),
		PublicKeyPEM:  string(publicPEM),
		CreatedAt:     time.Now().UTC(),
	}, nil
}

func parseRSAPrivateKey(value []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(value)
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	if privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return privateKey, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("DKIM private key is not RSA")
	}
	return privateKey, nil
}

func marshalPublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	encoded, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded}), nil
}

func identityFromKey(key entity.DKIMKey) infrastructureservices.DKIMIdentity {
	return infrastructureservices.DKIMIdentity{
		Domain:    key.Domain,
		Selector:  key.Selector,
		PublicKey: publicKeyValue(key.PublicKeyPEM),
	}
}

func publicKeyValue(value string) string {
	var builder strings.Builder
	for _, line := range strings.Split(value, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "-----") {
			continue
		}
		builder.WriteString(line)
	}
	return builder.String()
}
