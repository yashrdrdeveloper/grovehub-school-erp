package controller

import "net/http"

type WebhookController interface{GmailPush(http.ResponseWriter,*http.Request)}
