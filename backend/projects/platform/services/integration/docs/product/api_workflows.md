# Integration — API Execution Walkthroughs

**Preparation:** Run the backend with migrations, PostgreSQL, Infrastructure messaging and configured provider OAuth credentials. Set `BASE_URL` to an isolated backend address. Obtain an actually authenticated request context through supported middleware; do not spoof user or tenant headers. Examples intentionally omit production bearer tokens.

## 1. Discover providers

```bash
curl -i "$BASE_URL/integrations/providers"
curl -i "$BASE_URL/integrations/providers/GMAIL"
```

Expected: **200** with provider descriptors, available actions/capabilities and OAuth readiness. An unknown provider returns a documented error (typically **400**). Discovery does not create a connection or validate ownership.

## 2. Create pending provider connection

```bash
curl -i -X POST "$BASE_URL/integrations/connections" \
  -H "Content-Type: application/json" \
  -d '{"provider":"GMAIL","ownerGid":"<verified-owner-uuid>"}'
```

Expected **201** and a `PENDING_SETUP` connection with server-generated ID. `ownerGid` must match the **verified** actor context. An unauthorized request is denied (**403**), even if a syntactically valid owner ID was supplied. Capture `id` as `CONNECTION_ID`.

## 3. List and get connections

```bash
curl -i "$BASE_URL/integrations/connections"
curl -i "$BASE_URL/integrations/connections/$CONNECTION_ID"
```

Expected **200** with only the current authorized owner's connections and their connection/test state. Unowned connections must not be disclosed.

## 4. Start OAuth

```bash
curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/oauth/start"
```

Expected **200** with an authorization URL. Open it in a browser session using your *test provider account*; grant only configured scopes. Do not fabricate or reuse the returned one-time state.

## 5. Complete OAuth callback

The provider returns to the configured callback after authorization:

```text
GET /integrations/connections/oauth/{provider}/callback?state=<one-time-state>&code=<provider-code>
```

Expected **200** for successful code exchange and an active connection. Wrong/expired/replayed state, wrong audience or failed exchange must fail without activating the connection. **Never** print the authorization code or access token in your reproduction logs.

## 6. Run and inspect connection tests

```bash
curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/test" \
  -H "Content-Type: application/json" \
  -d '{"direction":"INBOUND"}'

curl -i "$BASE_URL/integrations/connections/$CONNECTION_ID/tests"
```

Expected **200** and a saved test result. Repeat with `{"direction":"OUTBOUND"}` only on a controlled test mailbox; it may send a self-addressed message. The history endpoint returns the latest results (up to 100 documented). A successful test is **not** proof that every future provider send is delivered.

## 7. Register/renew watch

```bash
curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/watch"
```

When Gmail watch support is configured, expect **200** and the subscription expiration metadata. If not configured, expect a documented provider-not-configured or unsupported result. Renewal should preserve authorized connection ownership.

## 8. Execute provider action

```bash
curl -i -X POST "$BASE_URL/integrations/connections/$CONNECTION_ID/actions/mail.profile" \
  -H "Content-Type: application/json" \
  -d '{"input":{}}'
```

Expected **200** on an active authorized connection with that action available. Use only names from provider discovery. Invalid input **400**, unowned resource **403**, inactive **409**, rate limited **429**, unsupported **501**, provider unavailable **503** as documented. See [mail_actions.md](mail_actions.md) for structured send/read examples.

## 9. Verified Gmail push webhook

`POST /integrations/webhooks/gmail` is called by **Google Pub/Sub** with a Google-signed push OIDC JWT and valid notification payload. A plain unsigned `curl` request is **not** a valid end-to-end test. Expected **204** after verified notification publication, **401** for invalid authentication and **503** when event publishing fails so upstream retries. This endpoint publishes a mailbox-change hint; it does **not** sync or store email messages.

## 10. Disconnect

```bash
curl -i -X DELETE "$BASE_URL/integrations/connections/$CONNECTION_ID"
```

Expected **200**, disconnected status and cleared stored OAuth credential. Subsequent owner actions on that connection must fail safely.

## 11. Acceptance checklist

Provider discovery, authorized create, list/get isolation, OAuth exchange/state validation, connection status, tested inbound/outbound access, watch expiration, provider action, verified webhook retry behavior and disconnect are independently reproducible. Compare actual response envelopes with OpenAPI; treat OAuth secrets and test messages as sensitive.
