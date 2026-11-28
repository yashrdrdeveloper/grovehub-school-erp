// Package parser discovers annotation occurrences without executing policies.
package parser

import (
	goparser "go/parser"
	"go/token"

	"grovehub/libs/core/annotation/internal/src/contracts"
	"grovehub/libs/core/annotation/internal/src/parser/helpers"
)

type Occurrence = contract.Occurrence

// ScanFile returns recognized annotations in source order.
func ScanFile(path string) ([]Occurrence, error) {
	set := token.NewFileSet()
	file, err := goparser.ParseFile(set, path, nil, goparser.ParseComments|goparser.AllErrors)
	if err != nil {
		return nil, err
	}

	targets := helpers.CollectTargets(file)

	return helpers.CollectOccurrences(set, file.Comments, targets)
}
