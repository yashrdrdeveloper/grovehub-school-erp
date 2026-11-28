package helpers

import "testing"

func TestRouteLimitForAbsentComment(t *testing.T) {
	limit, err := routeLimitFor(nil)
	if err != nil || limit != 0 {
		t.Fatalf("absent annotation: limit=%d err=%v", limit, err)
	}
}
