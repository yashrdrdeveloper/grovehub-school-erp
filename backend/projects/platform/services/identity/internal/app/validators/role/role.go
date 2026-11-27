package role

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/constants/persona"
	roledto "grovehub/services/identity/internal/app/dto/role"
)

var ErrInvalidDraft = errors.New("invalid role draft")
var ErrInvalidID = errors.New("invalid role identifier")

const MaxNameLength = 120
const MaxDescriptionLength = 1000

func NormalizeCreateDraft(request roledto.CreateDraft) (roledto.CreateDraft, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	if request.Name == "" || !utf8.ValidString(request.Name) || !utf8.ValidString(request.Description) || strings.ContainsRune(request.Name, 0) || strings.ContainsRune(request.Description, 0) {
		return roledto.CreateDraft{}, ErrInvalidDraft
	}
	if utf8.RuneCountInString(request.Name) > MaxNameLength || utf8.RuneCountInString(request.Description) > MaxDescriptionLength {
		return roledto.CreateDraft{}, ErrInvalidDraft
	}
	if len(request.Personas) == 0 || len(request.Personas) > len(persona.Catalog()) {
		return roledto.CreateDraft{}, ErrInvalidDraft
	}
	seen := make(map[persona.Code]bool, len(request.Personas))
	for _, code := range request.Personas {
		if !persona.Known(code) || seen[code] {
			return roledto.CreateDraft{}, ErrInvalidDraft
		}
		seen[code] = true
	}
	return request, nil
}

func NormalizeID(value string) (string, error) {
	identifier, err := uuid.Parse(value)
	if err != nil || identifier == uuid.Nil {
		return "", ErrInvalidID
	}
	return identifier.String(), nil
}

var ErrInvalidVersion = errors.New("invalid role version")

func NormalizeUpdateMetadata(request roledto.UpdateMetadata) (roledto.UpdateMetadata, error) {
	if request.Version < 1 {
		return roledto.UpdateMetadata{}, ErrInvalidVersion
	}
	draft, err := NormalizeCreateDraft(roledto.CreateDraft{Name: request.Name, Description: request.Description, Personas: request.Personas})
	if err != nil {
		return roledto.UpdateMetadata{}, err
	}
	request.Name = draft.Name
	request.Description = draft.Description
	request.Personas = draft.Personas
	return request, nil
}
