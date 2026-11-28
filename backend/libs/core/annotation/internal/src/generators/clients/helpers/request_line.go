package helpers

import (
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func requestLine(doc *ast.CommentGroup) (endpointMethod,bool,error) {
	if doc==nil { return endpointMethod{},false,nil }
	found:=false
	var result endpointMethod
	for _,comment := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if !strings.HasPrefix(line,"@RequestLine") { continue }
		if line!="@RequestLine" && !strings.HasPrefix(line,"@RequestLine ") { continue }
		if found { return result,true,fmt.Errorf("duplicate @RequestLine") }
		found=true
		args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(line,"@RequestLine")))
		if len(args)!=2 { return result,true,fmt.Errorf("@RequestLine requires METHOD /absolute/path") }
		result.verb,result.path=args[0],args[1]
		switch result.verb {
		case "GET","POST","PUT","PATCH","DELETE":
		default:
			return result,true,fmt.Errorf("unsupported @RequestLine verb %q",result.verb)
		}
		if !strings.HasPrefix(result.path,"/") || strings.Contains(result.path,"//") ||
			strings.ContainsAny(result.path,"?#") {
			return result,true,fmt.Errorf("invalid @RequestLine path %q",result.path)
		}
		for _,segment := range strings.Split(result.path,"/") {
			if segment=="."||segment==".." { return result,true,fmt.Errorf("unsafe @RequestLine path") }
		}
	}
	return result,found,nil
}

func endpointArguments(set *token.FileSet, fn *ast.FuncType, args []string, method endpointMethod,render layout.Renderer) (endpointMethod,error) {
	if len(args)!=2 || fn.Params==nil || len(fn.Params.List)!=2 || fn.Results==nil {
		return method,fmt.Errorf("@RequestLine requires context and one operation-specific request struct")
	}
	if !selectorType(fn.Params.List[0].Type,"context","Context") {
		return method,fmt.Errorf("first @RequestLine parameter must be context.Context")
	}
	// A request DTO is a named type (local or qualified). Its only permitted
	// exported fields are Params, Query and RequestBody; the transport enforces
	// their shape before dispatch, without any generic Empty placeholders.
	switch fn.Params.List[1].Type.(type) {
	case *ast.SelectorExpr, *ast.Ident:
	default:
		return method,fmt.Errorf("second @RequestLine parameter must be a named request struct")
	}
	results:=fn.Results.List
	if len(results)!=2 && len(results)!=3 {return method,fmt.Errorf("expected (response,error) or (response,bool,error)")}
	if !identType(results[len(results)-1].Type,"error") {return method,fmt.Errorf("last result must be error")}
	method.created=len(results)==3
	if method.created && !identType(results[1].Type,"bool") {
		return method,fmt.Errorf("second result must be bool (created flag)")
	}
	if method.created && method.verb!="POST" {return method,fmt.Errorf("created flag is supported only for POST")}
	if len(results[0].Names)>0 || len(results[len(results)-1].Names)>0 {
		return method,fmt.Errorf("named @RequestLine results are not supported")
	}
	result,err:=render(set,results[0].Type);if err!=nil{return method,err}
	method.resultType=result
	return method,nil
}
