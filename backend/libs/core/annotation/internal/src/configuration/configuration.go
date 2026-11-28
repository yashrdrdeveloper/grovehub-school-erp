// Package configuration provides typed startup-only environment loading.
package configuration

import (
	"fmt"
	"os"
	"reflect"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/configuration/helpers"
)

// Lookup can be injected without modifying the process environment.
type Lookup = helpers.Lookup

func FromEnv[T any]() (T, error) {
	return Decode[T](os.LookupEnv)
}

func Decode[T any](lookup Lookup) (T, error) {
	var value T

	if lookup == nil {
		return value, fmt.Errorf("configuration lookup is nil")
	}

	target := reflect.ValueOf(&value).Elem()
	if target.Kind() != reflect.Struct {
		return value, fmt.Errorf("configuration target must be a struct")
	}

	for index := 0; index < target.NumField(); index++ {
		field := target.Type().Field(index)

		if err := helpers.DecodeField(target.Field(index), field, lookup); err != nil {
			return value, err
		}
	}

	return value, nil
}
