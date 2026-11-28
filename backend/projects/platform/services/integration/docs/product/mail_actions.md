# Integration — Mail Action Walkthroughs

**Preconditions:** Active owner-authorized connection and supported action from provider catalog. Actual mail send can have external effects; use only approved test recipients and test/sandbox mailboxes.

## Generic action envelope

`POST /integrations/connections/{id}/actions/{action}` accepts `{"input":{...}}`. Do not include the access token, refresh token or OAuth state in the input; the connection owns encrypted credentials.

## Read profile, list and get

```bash
curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/actions/mail.profile" \
  -H "Content-Type: application/json" \
  -d '{"input":{}}'

curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/actions/mail.list" \
  -H "Content-Type: application/json" \
  -d '{"input":{"pageToken":""}}'

curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/actions/mail.get" \
  -H "Content-Type: application/json" \
  -d '{"input":{"messageId":"<real-provider-message-id>"}}'
```

Expected normalized provider output. Never invent provider-specific message IDs; get one from the real list result. Integration **does not persist** listed messages.

## Incremental change discovery

- Gmail: `mail.history` with `{"startHistoryId":"<provider-history-id>"}`.
- Outlook: `mail.delta` with `{"cursor":"<provider-cursor>"}`.

These are provider-specific read actions. The caller of Integration owns durable synchronization state; Integration does not commit remote cursors as the source of truth.

## Attachment retrieval

```bash
curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/actions/mail.attachment.get" \
  -H "Content-Type: application/json" \
  -d '{"input":{"messageId":"<real-id>","attachmentId":"<real-attachment-id>"}}'
```

Expected provider-normalized base64 content, subject to size limits, auth and provider availability.

## Send controlled test email

```bash
curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/actions/mail.send" \
  -H "Content-Type: application/json" \
  -d '{
    "input":{
      "to":["recipient@example.test"],
      "cc":[],
      "bcc":[],
      "subject":"Integration sandbox check",
      "body":{"text":"Test content from an isolated mailbox."}
    }
  }'
```

Replace the illustrative `example.test` address only with an **approved test inbox**. The API sends through the connected mailbox; an accepted provider send is **not a delivery receipt**. Optional HTML, reply-to and attachments require the OpenAPI structure; attachment bytes are **standard padded base64** and total limits apply. `from` is inferred from OAuth mailbox identity, not blindly supplied by a caller.

## Negative tests

Wrong/missing owner, disconnected credentials, unsupported action, invalid email, oversized/invalid attachment base64, missing body, expired access token, provider 429 and timeout should all return stable safe errors. Verify no token, raw MIME or provider internal response is exposed in logs.
