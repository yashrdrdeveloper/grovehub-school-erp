// Package clock provides stable timestamps shared by isolated unit tests.
package clock

import "time"

func Fixed() time.Time { return time.Date(2026,time.November,14,0,0,0,0,time.UTC) }
