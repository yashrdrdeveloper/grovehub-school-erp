package id

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

const maxUUIDv7UnixMilli = (1 << 48) - 1

var ErrUUIDv7TimestampOutOfRange = errors.New("uuidv7 timestamp is outside the supported Unix millisecond range")

func NewUUIDv7() (string, error) {
	return newUUIDv7(time.Now().UTC(), rand.Reader)
}

func newUUIDv7(now time.Time, random io.Reader) (string, error) {
	unixMilli := now.UnixMilli()
	if unixMilli < 0 || uint64(unixMilli) > maxUUIDv7UnixMilli {
		return "", ErrUUIDv7TimestampOutOfRange
	}

	var value [16]byte
	timestamp := uint64(unixMilli)
	value[0] = byte(timestamp >> 40)
	value[1] = byte(timestamp >> 32)
	value[2] = byte(timestamp >> 24)
	value[3] = byte(timestamp >> 16)
	value[4] = byte(timestamp >> 8)
	value[5] = byte(timestamp)

	if _, err := io.ReadFull(random, value[6:]); err != nil {
		return "", fmt.Errorf("generate uuidv7 randomness: %w", err)
	}

	value[6] = (value[6] & 0x0f) | 0x70
	value[8] = (value[8] & 0x3f) | 0x80

	return formatUUID(value), nil
}

func formatUUID(value [16]byte) string {
	var encoded [36]byte

	hex.Encode(encoded[0:8], value[0:4])
	encoded[8] = '-'
	hex.Encode(encoded[9:13], value[4:6])
	encoded[13] = '-'
	hex.Encode(encoded[14:18], value[6:8])
	encoded[18] = '-'
	hex.Encode(encoded[19:23], value[8:10])
	encoded[23] = '-'
	hex.Encode(encoded[24:36], value[10:16])

	return string(encoded[:])
}
