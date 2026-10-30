package impl

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	coreerrors "grovehub/libs/core/errors"
	infrastructureconfig "grovehub/services/infrastructure/config"
	infrastructureerror "grovehub/services/infrastructure/constants/error"
	"grovehub/services/infrastructure/entity"
	infrastructurerepository "grovehub/services/infrastructure/repository"
	infrastructureservices "grovehub/services/infrastructure/services"
	servicehelpers "grovehub/services/infrastructure/services/helpers"
)

const defaultDKIMBits = 2048

type DKIMService struct {
	repository infrastructurerepository.DKIMRepository
	config     infrastructureconfig.DKIM
}

func NewDKIMService(
	repository infrastructurerepository.DKIMRepository,
	config infrastructureconfig.DKIM,
) infrastructureservices.DKIMService {
	return &DKIMService{repository: repository, config: config}
}

func (service *DKIMService) Ensure(
	ctx context.Context,
	input infrastructureservices.EnsureDKIMInput,
) (infrastructureservices.DKIMIdentity, error) {
	domain := servicehelpers.NormalizeDKIMDomain(input.Domain)
	selector := servicehelpers.NormalizeDKIMSelector(input.Selector)
	if domain == "" || selector == "" {
		return infrastructureservices.DKIMIdentity{}, coreerrors.NewError(infrastructureerror.InvalidDKIMIdentity)
	}

	persisted, err := service.repository.FindByDomain(ctx, domain)
	if err == nil {
		return service.usePersisted(persisted, selector)
	}
	if !errors.Is(err, infrastructurerepository.ErrDKIMKeyNotFound) {
		return infrastructureservices.DKIMIdentity{}, coreerrors.Wrap(infrastructureerror.DKIMPrepareFailed, err)
	}

	key, found, err := service.readExistingFilesystemKey(domain, selector)
	if err != nil {
		var applicationError *coreerrors.Error
		if errors.As(err, &applicationError) {
			return infrastructureservices.DKIMIdentity{}, err
		}
		return infrastructureservices.DKIMIdentity{}, coreerrors.Wrap(infrastructureerror.DKIMPrepareFailed, err)
	}
	if !found {
		key, err = generateDKIMKey(domain, selector)
		if err != nil {
			return infrastructureservices.DKIMIdentity{}, coreerrors.Wrap(infrastructureerror.DKIMPrepareFailed, err)
		}
	}

	if err := service.repository.Create(ctx, key); err != nil {
		if !errors.Is(err, infrastructurerepository.ErrDKIMKeyAlreadyExists) {
			return infrastructureservices.DKIMIdentity{}, coreerrors.Wrap(infrastructureerror.DKIMPrepareFailed, err)
		}
		persisted, readErr := service.repository.FindByDomain(ctx, domain)
		if readErr != nil {
			return infrastructureservices.DKIMIdentity{}, coreerrors.Wrap(infrastructureerror.DKIMPrepareFailed, readErr)
		}
		return service.usePersisted(persisted, selector)
	}

	if err := service.materialize(key); err != nil {
		return infrastructureservices.DKIMIdentity{}, coreerrors.Wrap(infrastructureerror.DKIMPrepareFailed, err)
	}
	return identityFromKey(key), nil
}

func (service *DKIMService) usePersisted(
	key entity.DKIMKey,
	selector string,
) (infrastructureservices.DKIMIdentity, error) {
	if key.Selector != selector {
		return infrastructureservices.DKIMIdentity{}, coreerrors.NewError(infrastructureerror.DKIMSelectorConflict)
	}
	if err := service.materialize(key); err != nil {
		return infrastructureservices.DKIMIdentity{}, coreerrors.Wrap(infrastructureerror.DKIMPrepareFailed, err)
	}
	return identityFromKey(key), nil
}

func (service *DKIMService) readExistingFilesystemKey(
	domain string,
	selector string,
) (entity.DKIMKey, bool, error) {
	keyDir := filepath.Join(service.config.StateDir, domain)
	privatePEM, err := os.ReadFile(filepath.Join(keyDir, "private"))
	if errors.Is(err, os.ErrNotExist) {
		return entity.DKIMKey{}, false, nil
	}
	if err != nil {
		return entity.DKIMKey{}, false, err
	}
	if len(bytes.TrimSpace(privatePEM)) == 0 {
		return entity.DKIMKey{}, false, nil
	}

	selectorBytes, selectorErr := os.ReadFile(filepath.Join(keyDir, "selector"))
	if selectorErr == nil {
		existingSelector := servicehelpers.NormalizeDKIMSelector(string(selectorBytes))
		if existingSelector != "" && existingSelector != selector {
			return entity.DKIMKey{}, false, coreerrors.NewError(infrastructureerror.DKIMSelectorConflict)
		}
	} else if !errors.Is(selectorErr, os.ErrNotExist) {
		return entity.DKIMKey{}, false, selectorErr
	}

	privateKey, err := servicehelpers.ParseRSAPrivateKey(privatePEM)
	if err != nil {
		return entity.DKIMKey{}, false, fmt.Errorf("parse existing DKIM private key: %w", err)
	}
	publicPEM, err := servicehelpers.MarshalPublicKey(&privateKey.PublicKey)
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
	keyDir := filepath.Join(service.config.StateDir, key.Domain)
	if err := os.MkdirAll(keyDir, 0o700); err != nil {
		return err
	}
	privatePath := filepath.Join(keyDir, "private")
	if err := os.WriteFile(privatePath, []byte(key.PrivateKeyPEM), 0o600); err != nil {
		return err
	}
	if err := os.Chmod(privatePath, 0o600); err != nil {
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
	publicPEM, err := servicehelpers.MarshalPublicKey(&privateKey.PublicKey)
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

func identityFromKey(key entity.DKIMKey) infrastructureservices.DKIMIdentity {
	return infrastructureservices.DKIMIdentity{
		Domain:    key.Domain,
		Selector:  key.Selector,
		PublicKey: servicehelpers.PublicKeyValue(key.PublicKeyPEM),
	}
}
