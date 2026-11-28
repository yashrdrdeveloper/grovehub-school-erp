package helpers

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
)

func contains(all []string,target string)bool{
	for _,v:=range all {if v==target{return true}}
	return false
}

func compareFields(actual,expected []string,label string)error{
	sort.Strings(actual)
	sort.Strings(expected)
	if strings.Join(actual,",")!=strings.Join(expected,","){
		return fmt.Errorf("%s annotation must exactly match configured allowlist: configured=%v, declared=%v",label,actual,expected)
	}
	return nil
}

func validateConfig(dir,config string,filters,sorts []string,cursor string) error {
	lit,err:=findConfig(dir,config);if err!=nil{return err}
	scalars,err:=names(field(lit,"ScalarFields"));if err!=nil{return err}
	relations,err:=names(field(lit,"RelationFields"));if err!=nil{return err}
	if len(filters)>0 {
		if err:=compareFields(append(scalars,relations...),filters,"@Filterable");err!=nil{return err}
	}
	configuredSorts,err:=names(field(lit,"SortFields"));if err!=nil{return err}
	if len(sorts)>0 {
		if err:=compareFields(configuredSorts,sorts,"@Sortable");err!=nil{return err}
	}
	cursorExpr,err:=stringExpr(field(lit,"CursorExpression"));if err!=nil{return err}
	sortsNode,ok:=field(lit,"SortFields").(*ast.CompositeLit);if !ok{return fmt.Errorf("invalid SortFields")}
	for _,item:=range sortsNode.Elts {
		pair:=item.(*ast.KeyValueExpr)
		name,err:=stringExpr(pair.Key);if err!=nil{return err}
		if cursor!="" && name!=cursor {continue}
		value,ok:=pair.Value.(*ast.CompositeLit);if !ok{return fmt.Errorf("sort field %s must use a literal",cursor)}
		expression,err:=stringExpr(field(value,"Expression"));if err!=nil{return err}
		if expression==cursorExpr{return nil}
	}
	if cursor!="" {
		return fmt.Errorf("@Cursor sort %s must match CursorExpression",cursor)
	}
	return fmt.Errorf("no sortable field matches the stable CursorExpression %q",cursorExpr)
}
