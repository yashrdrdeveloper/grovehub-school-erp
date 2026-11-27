// Package providers defines provider-neutral discovery contracts.
package providers

import dto "grovehub/services/integration/internal/app/providers/dto"

// Provider is capability agnostic; concrete adapters own their action lists.
type Provider interface {
    Descriptor() dto.Descriptor
}
