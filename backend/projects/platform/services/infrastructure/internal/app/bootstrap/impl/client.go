package impl

import (
    client "grovehub/services/infrastructure/client/messaging"
    clientimpl "grovehub/services/infrastructure/internal/app/client/messaging/impl"
    messagingservice "grovehub/services/infrastructure/internal/app/services/messaging"
)

// Matches Foundation's resource-based client/<resource>/impl layout.
// Generated @Client proxy methods continue to live in the source overlay.
func (b *builder) newClient(service messagingservice.Service) client.Client {
    return clientimpl.NewMessagingClient(service)
}
