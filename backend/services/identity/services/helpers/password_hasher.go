package helpers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	passwordHashAlgorithm = "argon2id"
	passwordHashVersion   = argon2.Version
	passwordHashMemory    = 64 * 1024
	passwordHashTime      = 3
	passwordHashThreads   = 1
	passwordHashSaltLength = 16
	passwordHashKeyLength  = 32

	maxHashMemory    = 256 * 1024
	maxHashTime      = 10
	maxHashThreads   = 8
	maxHashKeyLength = 64
)

var ErrInvalidPasswordHash = errors.New("invalid password hash")

func HashPassword(password string) (string, error) {
	salt := make([]byte, passwordHashSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	derivedKey := argon2.IDKey(
		[]byte(password),
		salt,
		passwordHashTime,
		passwordHashMemory,
		passwordHashThreads,
		passwordHashKeyLength,
	)

	return fmt.Sprintf(
		"$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		passwordHashAlgorithm,
		passwordHashVersion,
		passwordHashMemory,
		passwordHashTime,
		passwordHashThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(derivedKey),
	), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parameters, salt, expectedHash, err := parsePasswordHash(encodedHash)
	if err != nil {
		return false, err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		parameters.time,
		parameters.memory,
		parameters.threads,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1, nil
}

type passwordHashParameters struct {
	memory  uint32
	time    uint32
	threads uint8
}

func parsePasswordHash(encodedHash string) (passwordHashParameters, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != passwordHashAlgorithm {
		return passwordHashParameters{}, nil, nil, ErrInvalidPasswordHash
	}

	version, err := parseHashUint(parts[2], "v=", 8)
	if err != nil || int(version) != passwordHashVersion {
		return passwordHashParameters{}, nil, nil, ErrInvalidPasswordHash
	}

	parameters, err := parseHashParameters(parts[3])
	if err != nil {
		return passwordHashParameters{}, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < passwordHashSaltLength {
		return passwordHashParameters{}, nil, nil, ErrInvalidPasswordHash
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expectedHash) < 1 || len(expectedHash) > maxHashKeyLength {
		return passwordHashParameters{}, nil, nil, ErrInvalidPasswordHash
	}

	return parameters, salt, expectedHash, nil
}

func parseHashParameters(value string) (passwordHashParameters, error) {
	parts := strings.Split(value, ",")
	if len(parts) != 3 {
		return passwordHashParameters{}, ErrInvalidPasswordHash
	}

	memory, err := parseHashUint(parts[0], "m=", 32)
	if err != nil || memory < 1 || memory > maxHashMemory {
		return passwordHashParameters{}, ErrInvalidPasswordHash
	}
	timeCost, err := parseHashUint(parts[1], "t=", 32)
	if err != nil || timeCost < 1 || timeCost > maxHashTime {
		return passwordHashParameters{}, ErrInvalidPasswordHash
	}
	threads, err := parseHashUint(parts[2], "p=", 8)
	if err != nil || threads < 1 || threads > maxHashThreads {
		return passwordHashParameters{}, ErrInvalidPasswordHash
	}

	return passwordHashParameters{
		memory:  uint32(memory),
		time:    uint32(timeCost),
		threads: uint8(threads),
	}, nil
}

func parseHashUint(value, prefix string, bitSize int) (uint64, error) {
	if !strings.HasPrefix(value, prefix) {
		return 0, ErrInvalidPasswordHash
	}
	number, err := strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, bitSize)
	if err != nil {
		return 0, ErrInvalidPasswordHash
	}
	return number, nil
}
