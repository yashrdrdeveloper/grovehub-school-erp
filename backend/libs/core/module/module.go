// Package module defines the shared HTTP-registration and shutdown contract
// for independently bootstrapped backend service modules.
package module

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/module/internal/src/contracts"

// Module is the minimal reusable lifecycle contract.
// Service-owned modules expose additional typed clients in their own interfaces.
type Module = contracts.Module
