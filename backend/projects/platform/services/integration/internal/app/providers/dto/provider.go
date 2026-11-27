package dto

// Action and Capability are generic identifiers; capability-specific constants
// (e.g. mail.send) belong under that capability's package.
type Action string

type Capability string

type Descriptor struct {
    ID string `json:"id"`
    DisplayName string `json:"displayName"`
    Capabilities []Capability `json:"capabilities"`
    AuthType string `json:"authType"`
}

type Session struct {
    ConnectionID string
    ProviderID string
    AccessToken string
}
