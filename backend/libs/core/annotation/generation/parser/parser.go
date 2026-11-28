// Package parser is a compatibility bridge for annotation source scanning.
package parser

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/parser"

type Occurrence = parser.Occurrence

func ScanFile(path string) ([]Occurrence, error) { return parser.ScanFile(path) }
