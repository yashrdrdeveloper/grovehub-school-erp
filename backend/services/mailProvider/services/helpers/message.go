package helpers

import (
	"strings"
	"time"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/mailProvider/entity"
)

type messageCursorValueExtractor func(entity.Message) string

var messageCursorExtractors = map[string]messageCursorValueExtractor{
	"messageId": func(message entity.Message) string { return string(message.GetMessageId()) },
	"gid":       func(message entity.Message) string { return string(message.GetGID()) },
	"senderId":  func(message entity.Message) string { return string(message.GetSenderId()) },
	"subject":   func(message entity.Message) string { return strings.ToLower(message.GetSubject()) },
	"status":    func(message entity.Message) string { return string(message.GetStatus()) },
	"createdAt": func(message entity.Message) string { return message.GetCreatedAt().UTC().Format(time.RFC3339Nano) },
	"updatedAt": func(message entity.Message) string { return message.GetUpdatedAt().UTC().Format(time.RFC3339Nano) },
}

func MessageCursorValue(message entity.Message, field string) string {
	extract, ok := messageCursorExtractors[field]
	if !ok {
		return ""
	}
	return extract(message)
}

func EncodeMessageSearchCursor(cursor coresearch.CursorPayload) (string, error) {
	return coresearch.EncodeCursor(cursor)
}
