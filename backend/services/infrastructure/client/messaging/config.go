package messaging

import config "grovehub/services/infrastructure/internal/app/config"

type Config = config.Messaging

const ProviderNATS = config.ProviderNATS

func DefaultMessaging() Config { return config.DefaultMessaging() }
