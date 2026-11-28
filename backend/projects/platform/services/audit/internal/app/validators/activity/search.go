package activity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/services/activity"
)

func validateID(value string, invalid error) error {
	id,err:=uuid.Parse(value)
	if err!=nil || id==uuid.Nil || id.String()!=value { return invalid }
	return nil
}

func normalizeID(_ coresearch.Operator,value string) (string,error) {
	return value,validateID(value,activity.ErrInvalidSearch)
}

func normalizeText(_ coresearch.Operator,value string) (string,error) {
	value=strings.TrimSpace(value)
	if value=="" || len(value)>128 { return "",activity.ErrInvalidSearch }
	return value,nil
}

func normalizeTime(_ coresearch.Operator,value string) (string,error) {
	at,err:=time.Parse(time.RFC3339Nano,value)
	if err!=nil { return "",activity.ErrInvalidSearch }
	return at.UTC().Format(time.RFC3339Nano),nil
}

func sortFields() map[string]coresearch.CursorValueValidator {
	return map[string]coresearch.CursorValueValidator{
		"occurredAt":coresearch.ValidateTimeCursor,"ingestedAt":coresearch.ValidateTimeCursor,
		"auditId":validateID,"assetType":coresearch.ValidateNonEmptyCursor,"operation":coresearch.ValidateNonEmptyCursor,"source":coresearch.ValidateNonEmptyCursor,"auditType":coresearch.ValidateNonEmptyCursor,
	}
}

func filterFields() map[string]coresearch.FieldSpec {
	fields:=map[string]coresearch.FieldSpec{}
	for _,field:=range []string{"auditId","assetId","actorId","personId","subjectAccountId","occurrenceId"} {
		fields[field]=coresearch.FieldSpec{Operators:coresearch.EqualityOperators(),Normalizer:normalizeID}
	}
	for _,field:=range []string{"assetType","auditType","source","operation"} {
		fields[field]=coresearch.FieldSpec{Operators:coresearch.TextOperators(),Normalizer:normalizeText}
	}
	for _,field:=range []string{"occurredAt","ingestedAt"} {
		fields[field]=coresearch.FieldSpec{Operators:coresearch.TimeOperators(),Normalizer:normalizeTime}
	}
	return fields
}

func Normalize(input dto.Search) (coresearch.Query[string],error) {
	invalid:=activity.ErrInvalidSearch
	if input.Scope!="" && input.Scope!="ME" { return coresearch.Query[string]{},invalid }
	terms,err:=coresearch.NormalizeSearchTerms(input.SearchText,256,10,invalid)
	if err!=nil { return coresearch.Query[string]{},err }
	page,rows,offset,err:=coresearch.ResolvePage(input.PageInfo,invalid)
	if err!=nil || offset>100000 { return coresearch.Query[string]{},invalid }
	fields:=sortFields()
	sorts,err:=coresearch.NormalizeSorts(input.Sort,coresearch.QuerySort{Field:"occurredAt",Direction:coresearch.DirectionDescending},fields,5,invalid)
	if err!=nil { return coresearch.Query[string]{},err }
	filters,err:=coresearch.NormalizeFilters(input.Filters,filterFields(),20,50,invalid)
	if err!=nil { return coresearch.Query[string]{},err }
	cursor,err:=coresearch.NormalizeMultiCursor(input.PageInfo.Cursor,sorts,fields,validateID,invalid)
	if err!=nil { return coresearch.Query[string]{},err }
	if err:=coresearch.ValidatePageCursor(input.PageInfo,page,rows,invalid);err!=nil { return coresearch.Query[string]{},err }
	if input.PageInfo.PageNumber==0 { page=0 }
	return coresearch.Query[string]{SearchTerms:terms,Filters:filters,Sort:sorts[0],SecondarySorts:sorts[1:],Cursor:cursor,Limit:rows,PageNumber:page,Offset:offset},nil
}
