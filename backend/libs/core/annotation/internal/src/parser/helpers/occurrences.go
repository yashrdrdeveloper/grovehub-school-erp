package helpers

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/registry"
)

// CollectOccurrences preserves source order and reports invalid directives
// with their precise source location.
func CollectOccurrences(set *token.FileSet, groups []*ast.CommentGroup, targets Targets) ([]contract.Occurrence, error) {
	var result []contract.Occurrence

	for _, group := range groups {
		occurrences, err := collectGroup(set, group, targets.For(group))
		if err != nil {
			return nil, err
		}

		result = append(result, occurrences...)
	}

	return result, nil
}

func collectGroup(set *token.FileSet, group *ast.CommentGroup, target contract.Target) ([]contract.Occurrence, error) {
	var result []contract.Occurrence

	for _, comment := range group.List {
		directive, found, err := registry.ParseComment(comment.Text)
		position := set.Position(comment.Pos())
		if err != nil {
			return nil, fmt.Errorf("%s:%d:%d: %w", position.Filename, position.Line, position.Column, err)
		}
		if !found {
			continue
		}

		result = append(result, contract.Occurrence{
			Directive: directive,
			Target:    target,
			File:      position.Filename,
			Line:      position.Line,
			Column:    position.Column,
		})
	}

	return result, nil
}
