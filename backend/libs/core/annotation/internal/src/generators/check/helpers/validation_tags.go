package helpers

import (
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"grovehub/libs/core/annotation/internal/src/contracts"
	"grovehub/libs/core/annotation/internal/src/registry"
)

func CheckValidationTags(path string) error {
	source, err := layout.ParseSource(path,goparser.ParseComments)
	if err != nil {
		return err
	}

	var validationError error

	ast.Inspect(source.Tree, func(node ast.Node) bool {
		if validationError != nil {
			return false
		}

		field, ok := node.(*ast.Field)
		if !ok || field.Doc == nil {
			return true
		}

		validationError = validateFieldTags(path, field, source.Tags[field])

		return validationError == nil
	})

	return validationError
}

func validateFieldTags(path string, field *ast.Field,raw string) error {
	tags := strings.Split(reflect.StructTag(raw).Get("validate"),",")

	for _, comment := range field.Doc.List {
		directive, found, err := registry.ParseComment(comment.Text)
		if err != nil {
			return err
		}
		if !found {
			continue
		}

		rule, err := validationRule(path, directive)
		if err != nil {
			return err
		}
		if rule == "" {
			continue
		}
		if !slices.Contains(tags, rule) {
			return fmt.Errorf("%s: @%s must have validate:%s", path, directive.Name, rule)
		}
	}

	return nil
}

func validationRule(path string, directive contract.Directive) (string, error) {
	switch directive.Name {
	case "Valid":
		return "valid", nil
	case "NotBlank":
		return "notblank", nil
	case "NotNull", "NonNull":
		return "nonnull", nil
	case "Size":
		return sizeRule(path, directive.Arguments)
	default:
		return "", nil
	}
}

func sizeRule(path, arguments string) (string, error) {
	parts := strings.Fields(arguments)
	if len(parts) != 1 || !strings.HasPrefix(parts[0], "max=") {
		return "", fmt.Errorf("%s: @Size requires max=N", path)
	}

	size, err := strconv.Atoi(strings.TrimPrefix(parts[0], "max="))
	if err != nil || size < 0 {
		return "", fmt.Errorf("%s: invalid @Size maximum", path)
	}

	return fmt.Sprintf("maxrunes=%d", size), nil
}
