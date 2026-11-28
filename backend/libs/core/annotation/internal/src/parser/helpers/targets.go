package helpers

import (
	"go/ast"
	"go/token"

	"grovehub/libs/core/annotation/internal/src/contracts"
)

type Targets map[*ast.CommentGroup]contract.Target

// CollectTargets visits the AST once and records annotation comment ownership.
func CollectTargets(file *ast.File) Targets {
	targets := make(Targets)

	ast.Inspect(file, func(node ast.Node) bool {
		switch decl := node.(type) {
		case *ast.File:
			targets.record(decl.Doc, contract.PackageTarget)
		case *ast.FuncDecl:
			targets.record(decl.Doc, functionTarget(decl))
		case *ast.GenDecl:
			targets.record(decl.Doc, declarationTarget(decl))
		case *ast.TypeSpec:
			targets.record(decl.Doc, typeTarget(decl.Type))
		case *ast.InterfaceType:
			targets.recordInterfaceMethods(decl)
		case *ast.Field:
			targets.recordField(decl)
		}

		return true
	})

	return targets
}

func (targets Targets) record(group *ast.CommentGroup, target contract.Target) {
	if group == nil {
		return
	}

	targets[group] = target
}

func (targets Targets) For(group *ast.CommentGroup) contract.Target {
	if target, exists := targets[group]; exists {
		return target
	}

	return contract.FileTarget
}

func (targets Targets) recordField(field *ast.Field) {
	for _, group := range []*ast.CommentGroup{field.Doc, field.Comment} {
		if _, exists := targets[group]; exists {
			continue
		}

		targets.record(group, contract.FieldTarget)
	}
}

func (targets Targets) recordInterfaceMethods(decl *ast.InterfaceType) {
	for _, method := range decl.Methods.List {
		if !isNamedMethod(method) {
			continue
		}

		targets.record(method.Doc, contract.MethodTarget)
		targets.record(method.Comment, contract.MethodTarget)
	}
}

func isNamedMethod(field *ast.Field) bool {
	_, isFunction := field.Type.(*ast.FuncType)

	return len(field.Names) == 1 && isFunction
}

func functionTarget(decl *ast.FuncDecl) contract.Target {
	if decl.Recv != nil {
		return contract.MethodTarget
	}

	return contract.FunctionTarget
}

func declarationTarget(decl *ast.GenDecl) contract.Target {
	if decl.Tok != token.TYPE {
		return contract.FileTarget
	}

	if len(decl.Specs) != 1 {
		return contract.TypeTarget
	}

	spec, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return contract.TypeTarget
	}

	return typeTarget(spec.Type)
}

func typeTarget(expr ast.Expr) contract.Target {
	if _, ok := expr.(*ast.InterfaceType); ok {
		return contract.InterfaceTarget
	}

	return contract.TypeTarget
}
