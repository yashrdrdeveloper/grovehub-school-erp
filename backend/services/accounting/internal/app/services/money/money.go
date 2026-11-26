// Package money implements exact minor-unit arithmetic for the initial INR book.
package money

import (
 "errors"
 "math"
)

type Amount int64

var ErrOverflow = errors.New("monetary amount overflow")

func Add(a, b Amount) (Amount, error) {
 if (b > 0 && a > Amount(math.MaxInt64)-b) || (b < 0 && a < Amount(math.MinInt64)-b) {
  return 0, ErrOverflow
 }
 return a+b, nil
}

func Subtract(a, b Amount) (Amount, error) {
 if (b > 0 && a < Amount(math.MinInt64)+b) || (b < 0 && a > Amount(math.MaxInt64)+b) {
  return 0, ErrOverflow
 }
 return a-b, nil
}
