package helpers

import (
	"strings"
	"time"

	coreerrors "grovehub/libs/core/errors"
	corenormalize "grovehub/libs/core/normalize"
	coresearch "grovehub/libs/core/search"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
)

const (
	maxRecipients              = 100
	maxMessageSearchTextLength = 256
	maxMessageSearchTerms      = 10
	maxMessageSearchFilters    = 20
	maxMessageFilterValues     = 50
)

type messageValueNormalizer func(coresearch.Operator, string) (string, error)
type messageCursorValueValidator func(string) error

type messageSearchFieldSpec struct {
	operators  map[coresearch.Operator]struct{}
	normalizer messageValueNormalizer
}

type messageSortFieldSpec struct {
	validate messageCursorValueValidator
}

var messageSearchFields = map[string]messageSearchFieldSpec{
	"messageId": {
		operators: messageOperatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
		),
		normalizer: normalizeMessageUUIDSearchValue,
	},
	"gid": {
		operators: messageOperatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
		),
		normalizer: normalizeMessageUUIDSearchValue,
	},
	"senderId": {
		operators: messageOperatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
		),
		normalizer: normalizeMessageUUIDSearchValue,
	},
	"subject": {
		operators: messageOperatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
			coresearch.OperatorContains,
			coresearch.OperatorStartsWith,
		),
		normalizer: normalizeMessageTextSearchValue,
	},
	"status": {
		operators: messageOperatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
		),
		normalizer: normalizeMessageStatusSearchValue,
	},
	"createdAt": {
		operators: messageOperatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorGreaterThan,
			coresearch.OperatorGTE,
			coresearch.OperatorLessThan,
			coresearch.OperatorLTE,
		),
		normalizer: normalizeMessageTimeSearchValue,
	},
	"updatedAt": {
		operators: messageOperatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorGreaterThan,
			coresearch.OperatorGTE,
			coresearch.OperatorLessThan,
			coresearch.OperatorLTE,
		),
		normalizer: normalizeMessageTimeSearchValue,
	},
}

var messageSortFields = map[string]messageSortFieldSpec{
	"messageId": {validate: validateMessageUUIDCursorValue},
	"gid":       {validate: validateMessageUUIDCursorValue},
	"senderId":  {validate: validateMessageUUIDCursorValue},
	"subject":   {validate: validateMessageTextCursorValue},
	"status":    {validate: validateMessageStatusCursorValue},
	"createdAt": {validate: validateMessageTimeCursorValue},
	"updatedAt": {validate: validateMessageTimeCursorValue},
}

var messageValueCountValidators = map[coresearch.Operator]func(int) bool{
	coresearch.OperatorIn:    func(count int) bool { return count > 0 },
	coresearch.OperatorNotIn: func(count int) bool { return count > 0 },
}

func NormalizeGID(value string) (entity.GID, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	return entity.GID(normalized), nil
}

func NormalizeMessageId(value string) (entity.MessageId, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	return entity.MessageId(normalized), nil
}

func NormalizeRecipients(values []string) ([]string, error) {
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		email, err := corenormalize.Email(value)
		if err != nil {
			return nil, coreerrors.NewError(mailprovidererror.InvalidMessage)
		}
		if _, ok := seen[email]; ok {
			continue
		}
		seen[email] = struct{}{}
		normalized = append(normalized, email)
	}
	return normalized, nil
}

func ValidateMessageContent(to, cc, bcc []string, subject, textBody, htmlBody string) error {
	if len(to) == 0 || len(to)+len(cc)+len(bcc) > maxRecipients {
		return coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	if strings.TrimSpace(subject) == "" || len([]rune(subject)) > 998 {
		return coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	if strings.TrimSpace(textBody) == "" && strings.TrimSpace(htmlBody) == "" {
		return coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	return nil
}

func NormalizeMessageSearchRequest(
	request coresearch.Request,
) (coresearch.Query[entity.MessageId], coresearch.CursorPayload, error) {
	searchText := strings.TrimSpace(request.SearchText)
	if len(searchText) > maxMessageSearchTextLength {
		return coresearch.Query[entity.MessageId]{}, coresearch.CursorPayload{}, invalidMessageSearch()
	}
	searchTerms := strings.Fields(strings.ToLower(searchText))
	if len(searchTerms) > maxMessageSearchTerms || len(request.Filters) > maxMessageSearchFilters {
		return coresearch.Query[entity.MessageId]{}, coresearch.CursorPayload{}, invalidMessageSearch()
	}

	rows := request.PageInfo.Rows
	if rows == 0 {
		rows = coresearch.DefaultRows
	}
	if rows < 1 || rows > coresearch.MaxRows {
		return coresearch.Query[entity.MessageId]{}, coresearch.CursorPayload{}, invalidMessageSearch()
	}

	sort, err := normalizeMessageSort(request.Sort)
	if err != nil {
		return coresearch.Query[entity.MessageId]{}, coresearch.CursorPayload{}, err
	}
	filters, err := normalizeMessageFilters(request.Filters)
	if err != nil {
		return coresearch.Query[entity.MessageId]{}, coresearch.CursorPayload{}, err
	}
	cursor, payload, err := normalizeMessageCursor(request.PageInfo.Cursor, sort)
	if err != nil {
		return coresearch.Query[entity.MessageId]{}, coresearch.CursorPayload{}, err
	}

	return coresearch.Query[entity.MessageId]{
		SearchTerms: searchTerms,
		Filters:     filters,
		Sort:        sort,
		Cursor:      cursor,
		Limit:       rows,
	}, payload, nil
}

func normalizeMessageSort(sortRequest []coresearch.Sort) (coresearch.QuerySort, error) {
	sort := coresearch.QuerySort{Field: "createdAt", Direction: coresearch.DirectionDescending}
	if len(sortRequest) > 1 {
		return coresearch.QuerySort{}, invalidMessageSearch()
	}
	if len(sortRequest) == 1 {
		sort.Field = strings.TrimSpace(sortRequest[0].Field)
		sort.Direction = sortRequest[0].Direction
		if sort.Direction == "" {
			sort.Direction = coresearch.DirectionDescending
		}
	}
	if _, ok := messageSortFields[sort.Field]; !ok ||
		(sort.Direction != coresearch.DirectionAscending && sort.Direction != coresearch.DirectionDescending) {
		return coresearch.QuerySort{}, invalidMessageSearch()
	}
	return sort, nil
}

func normalizeMessageFilters(filters []coresearch.Filter) ([]coresearch.QueryFilter, error) {
	normalized := make([]coresearch.QueryFilter, 0, len(filters))
	for _, filter := range filters {
		field := strings.TrimSpace(filter.Field)
		spec, ok := messageSearchFields[field]
		if !ok {
			return nil, invalidMessageSearch()
		}
		if _, ok := spec.operators[filter.Operator]; !ok ||
			!validMessageSearchValueCount(filter.Operator, len(filter.Values)) ||
			len(filter.Values) > maxMessageFilterValues {
			return nil, invalidMessageSearch()
		}

		values := make([]string, 0, len(filter.Values))
		for _, value := range filter.Values {
			normalizedValue, err := spec.normalizer(filter.Operator, value)
			if err != nil {
				return nil, err
			}
			values = append(values, normalizedValue)
		}
		normalized = append(normalized, coresearch.QueryFilter{
			Field:    field,
			Operator: filter.Operator,
			Values:   values,
		})
	}
	return normalized, nil
}

func normalizeMessageCursor(
	value string,
	sort coresearch.QuerySort,
) (*coresearch.Cursor[entity.MessageId], coresearch.CursorPayload, error) {
	payload := coresearch.CursorPayload{Field: sort.Field, Direction: sort.Direction}
	if strings.TrimSpace(value) == "" {
		return nil, payload, nil
	}

	decoded, err := coresearch.DecodeCursor(value)
	if err != nil || decoded.Field != sort.Field || decoded.Direction != sort.Direction {
		return nil, coresearch.CursorPayload{}, invalidMessageSearch()
	}
	spec := messageSortFields[sort.Field]
	if err := spec.validate(decoded.Value); err != nil || validateMessageUUIDCursorValue(decoded.CursorId) != nil {
		return nil, coresearch.CursorPayload{}, invalidMessageSearch()
	}

	return &coresearch.Cursor[entity.MessageId]{
		Value:    decoded.Value,
		CursorId: entity.MessageId(decoded.CursorId),
	}, decoded, nil
}

func validMessageSearchValueCount(operator coresearch.Operator, count int) bool {
	if validate, ok := messageValueCountValidators[operator]; ok {
		return validate(count)
	}
	return count == 1
}

func messageOperatorSet(operators ...coresearch.Operator) map[coresearch.Operator]struct{} {
	set := make(map[coresearch.Operator]struct{}, len(operators))
	for _, operator := range operators {
		set[operator] = struct{}{}
	}
	return set
}

func normalizeMessageUUIDSearchValue(_ coresearch.Operator, value string) (string, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", invalidMessageSearch()
	}
	return normalized, nil
}

func normalizeMessageTextSearchValue(_ coresearch.Operator, value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", invalidMessageSearch()
	}
	return normalized, nil
}

func normalizeMessageStatusSearchValue(_ coresearch.Operator, value string) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	switch entity.MessageStatus(status) {
	case entity.MessageStatusQueued, entity.MessageStatusSubmitted, entity.MessageStatusFailed:
		return status, nil
	default:
		return "", invalidMessageSearch()
	}
}

func normalizeMessageTimeSearchValue(_ coresearch.Operator, value string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return "", invalidMessageSearch()
	}
	return parsed.UTC().Format(time.RFC3339Nano), nil
}

func validateMessageUUIDCursorValue(value string) error {
	_, err := corenormalize.UUIDv7(value)
	return err
}

func validateMessageTextCursorValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return invalidMessageSearch()
	}
	return nil
}

func validateMessageStatusCursorValue(value string) error {
	_, err := normalizeMessageStatusSearchValue(coresearch.OperatorEquals, value)
	return err
}

func validateMessageTimeCursorValue(value string) error {
	_, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(value))
	return err
}

func invalidMessageSearch() error {
	return coreerrors.NewError(mailprovidererror.InvalidMessageSearch)
}
