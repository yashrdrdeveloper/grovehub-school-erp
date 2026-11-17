package impl

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// IntegrationMailSender invokes the existing mail.send provider action through
// a dedicated system-owned and owner-authorized Integration connection.
// It does not accept an end-user OAuth credential or arbitrary sender identity.
type IntegrationMailSender struct {
    client *http.Client
    endpoint string
    serviceToken string
}

func NewIntegrationMailSender(rawBaseURL string, connectionID string, serviceToken string) (*IntegrationMailSender, error) {
    base, err := url.Parse(strings.TrimRight(rawBaseURL, "/"))
    if err != nil || base.Scheme != "https" || base.Host == "" || connectionID == "" || serviceToken == "" {
        return nil, errors.New("system mail requires HTTPS, an authorized connection and a service token")
    }
    endpoint := strings.TrimRight(base.String(), "/") + "/integrations/connections/" +
        url.PathEscape(connectionID) + "/actions/mail.send"
    return &IntegrationMailSender{
        client: &http.Client{Timeout: 15*time.Second},
        endpoint: endpoint, serviceToken: serviceToken,
    }, nil
}

func (sender *IntegrationMailSender) SendPasswordSetup(ctx context.Context, recipient string, setupLink string) error {
    if recipient == "" || setupLink == "" { return errors.New("invalid password invitation") }
    payload := map[string]any{
        "input": map[string]any{
            "to": []string{recipient},
            "cc": []string{},
            "bcc": []string{},
            "subject": "Set up your GrooveHub password",
            "body": map[string]string{
                "text": "An administrator created an account for you. Set your password using this one-time link (valid for 30 minutes): " + setupLink +
                    "\nIf you did not expect this message, you may ignore it.",
            },
        },
    }
    encoded, err := json.Marshal(payload)
    if err != nil { return err }
    request, err := http.NewRequestWithContext(ctx, http.MethodPost, sender.endpoint, bytes.NewReader(encoded))
    if err != nil { return err }
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Authorization", "Bearer "+sender.serviceToken)
    response, err := sender.client.Do(request)
    if err != nil { return errors.New("system mail delivery request failed") }
    defer response.Body.Close()
    if response.StatusCode < 200 || response.StatusCode >= 300 {
        // Do not log or expose provider response bodies, mail tokens or secrets.
        return fmt.Errorf("system mail provider rejected request (status %d)", response.StatusCode)
    }
    return nil
}
