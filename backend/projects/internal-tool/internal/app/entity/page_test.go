package entity

import "testing"

func TestPage(t *testing.T){
 p,ok:=(PageRequest{}).Normalize();if !ok||p.PageNumber!=1||p.RowSize!=25{t.Fatal("defaults",p)}
 for _,bad:=range []PageRequest{{PageNumber:-1},{RowSize:101},{PageNumber:10001}}{if _,ok:=bad.Normalize();ok{t.Fatal("accepted",bad)}}
 p=PageRequest{PageNumber:2,RowSize:2}
 if p.Offset()!=2{t.Fatal("offset")}
 page:=NewPage(p,[]int{1,2,3});if !page.HasMore||len(page.Records)!=2{t.Fatal("page",page)}
}
