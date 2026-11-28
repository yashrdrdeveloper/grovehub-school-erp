// Package gateway owns application-neutral route tiers and workload verification.
package gateway

import (
	"context"
	"net/http"
)

type Tier string

const (
	Public     Tier = "public"
	Protected  Tier = "protected"
	Restricted Tier = "restricted"
)

type WorkloadVerifier interface {
	VerifyWorkload(ctx context.Context, request *http.Request) (context.Context, error)
}
