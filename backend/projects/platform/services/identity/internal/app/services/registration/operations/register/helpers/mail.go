package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"

	integrationclient "grovehub/services/integration/client/connection"
	"grovehub/services/identity/internal/resources"
)

type Mail struct {
	DeliveryID string
	Title string
	Name string
	Description string
	Button string
	URL string
	Minutes int
}

func RenderMail(mail Mail) (string, error) {
	template, err := template.ParseFS(resources.Templates, "templates/registration/password_setup.html")
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := template.Execute(&buffer, mail); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// Send renders the resource-owned template and delegates to Integration's
// public typed client. The underlying provider owns Gmail/Outlook transport.
func Send(ctx context.Context, client integrationclient.Client, tenantID string, recipient string, mail Mail) error {
	if client == nil || tenantID == "" {
		return errors.New("system mail service or tenant is unavailable")
	}

	html, err := RenderMail(mail)
	if err != nil {
		return err
	}

	input := map[string]any{
		"to": []string{recipient},
		"subject": mail.Title,
		"body": map[string]string{
			"html": html,
			"text": mail.Description + "\n" + mail.URL,
		},
	}
	if mail.DeliveryID != "" {
		input["deliveryId"] = mail.DeliveryID
	}
	payload, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = client.ExecuteSystemMail(ctx, tenantID, json.RawMessage(payload))
	return err
}
