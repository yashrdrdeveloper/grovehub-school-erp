package impl

import client "grovehub/services/identity/client/invitation"

type invitationClient struct {
	reader client.Client
}

func NewInvitationClient(reader client.Client) client.Client {
	return &invitationClient{reader: reader}
}
