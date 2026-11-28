package operator_test

import (
 "testing"
 dto "grovehub/services/audit/internal/app/dto/operator"
 validator "grovehub/services/audit/internal/app/validators/operator"
)
func TestReceiptFilters(t *testing.T){
 base:=dto.ReceiptSearchInput{SearchInput:dto.SearchInput{PageNumber:1,RowSize:25}}
 if !validator.ValidReceiptSearch(base){t.Fatal("global search denied")}
 base.SchoolID="e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
 for _,op:=range []string{"","RENAME","ACTIVATE","DISABLE"}{base.Operation=op;if !validator.ValidReceiptSearch(base){t.Fatal(op)}}
 base.Operation="DELETE";if validator.ValidReceiptSearch(base){t.Fatal("unknown operation")}
 base.Operation="";base.SchoolID="00000000-0000-0000-0000-000000000000";if validator.ValidReceiptSearch(base){t.Fatal("nil UUID")}
 base.SchoolID="E93696CC-15F0-4E06-9577-0D5CFC2A83C8";if validator.ValidReceiptSearch(base){t.Fatal("noncanonical UUID")}
 base.SchoolID="";base.PageNumber=0;if validator.ValidReceiptSearch(base){t.Fatal("invalid pagination")}
}
