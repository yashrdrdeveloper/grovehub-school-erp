package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

const MaxInputBytes = 16 * 1024 * 1024

var ErrInvalidInput = errors.New("invalid migration input or reviewed SHA-256 mismatch")

func ReadReviewed(path string, expected string) ([]byte, error) {
	if len(expected) != 64 || strings.ToLower(expected) != expected {
		return nil, ErrInvalidInput
	}
	if _, err := hex.DecodeString(expected); err != nil {
		return nil, ErrInvalidInput
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, MaxInputBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > MaxInputBytes || Digest(data) != expected {
		return nil, ErrInvalidInput
	}
	return data, nil
}

func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func Decode(data []byte, result any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(result); err != nil {
		return ErrInvalidInput
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return ErrInvalidInput
	}
	return nil
}
