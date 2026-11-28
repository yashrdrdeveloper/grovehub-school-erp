package audit
import("testing";"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity")
func TestReceiptInput(t *testing.T){
 for _,tc:=range []struct{id,op string;valid bool}{
 {"","",true},{"e93696cc-15f0-4e06-9577-0d5cfc2a83c8","DISABLE",true},{"invalid","",false},{"00000000-0000-0000-0000-000000000000","",false},{"","DELETE",false},
 }{in,ok:=(ReceiptSearchInput{PageRequest:entity.PageRequest{},SchoolID:tc.id,Operation:tc.op}).Normalize();if ok!=tc.valid||in.PageNumber!=1||in.RowSize!=25{t.Fatal(tc,in,ok)}}
}
