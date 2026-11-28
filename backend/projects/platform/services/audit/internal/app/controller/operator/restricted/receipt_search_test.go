package restricted
import "testing"
func TestReceiptFilters(t *testing.T){
 base:=ReceiptSearchInput{SearchInput:SearchInput{PageNumber:1,RowSize:25}}
 if !base.Valid(){t.Fatal("global search denied")}
 base.SchoolID="e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
 for _,op:=range []string{"","RENAME","ACTIVATE","DISABLE"}{base.Operation=op;if !base.Valid(){t.Fatal(op)}}
 base.Operation="DELETE";if base.Valid(){t.Fatal("unknown operation")}
 base.Operation="";base.SchoolID="00000000-0000-0000-0000-000000000000";if base.Valid(){t.Fatal("nil UUID")}
 base.SchoolID="E93696CC-15F0-4E06-9577-0D5CFC2A83C8";if base.Valid(){t.Fatal("noncanonical UUID")}
 base.SchoolID="";base.PageNumber=0;if base.Valid(){t.Fatal("invalid pagination")}
}
