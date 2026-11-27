package helpers

import (
    maildto "grovehub/services/integration/internal/app/providers/mailprovider/dto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"

)

type TokenCipher struct{ key []byte }

func NewTokenCipher(key []byte) (TokenCipher, error) {
	if len(key) != 32 {
		return TokenCipher{}, errors.New("integration encryption key must have 32 bytes")
	}
	return TokenCipher{key: append([]byte(nil), key...)}, nil
}

func (c TokenCipher) Seal(token maildto.OAuthToken) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	payload, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, payload, nil), nil
}

func (c TokenCipher) Open(ciphertext []byte) (maildto.OAuthToken, error) {
	var token maildto.OAuthToken
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return token, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return token, errors.New("invalid credential ciphertext")
	}
	payload, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return token, err
	}
	err = json.Unmarshal(payload, &token)
	return token, err
}
