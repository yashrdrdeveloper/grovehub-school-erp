package helpers

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func assignValue(target reflect.Value, field, key, text string) error {
	switch {
	case target.Kind() == reflect.String:
		target.SetString(text)
		return nil
	case target.Type() == reflect.TypeOf(time.Duration(0)):
		return assignDuration(target, key, text)
	case target.Kind() == reflect.Bool:
		return assignBoolean(target, key, text)
	case target.Kind() >= reflect.Int && target.Kind() <= reflect.Int64:
		return assignInteger(target, key, text)
	default:
		return fmt.Errorf("unsupported configuration type for %s", field)
	}
}

func assignDuration(target reflect.Value, key, text string) error {
	value, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("%s: invalid duration", key)
	}

	target.SetInt(int64(value))
	return nil
}

func assignBoolean(target reflect.Value, key, text string) error {
	value, err := strconv.ParseBool(text)
	if err != nil {
		return fmt.Errorf("%s: invalid boolean", key)
	}

	target.SetBool(value)
	return nil
}

func assignInteger(target reflect.Value, key, text string) error {
	value, err := strconv.ParseInt(text, 10, target.Type().Bits())
	if err != nil {
		return fmt.Errorf("%s: invalid integer", key)
	}

	target.SetInt(value)
	return nil
}
