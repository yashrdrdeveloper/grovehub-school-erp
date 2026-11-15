package messaging

import inframessaging "grovehub/libs/infrastructure/messaging"

// Internal configuration uses the Infrastructure-owned public contract.
type Messaging = inframessaging.Config

const ProviderNATS = inframessaging.ProviderNATS

func DefaultMessaging() Messaging { return inframessaging.FromEnv() }
