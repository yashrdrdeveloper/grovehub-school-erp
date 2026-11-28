package impl_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"grovehub/services/identity/client/invitation"
	"grovehub/services/integration/client/mail"
	"grovehub/services/integration/internal/app/repository/maildelivery"
	"grovehub/services/integration/internal/app/services/mail/operations/reconciliation/impl"
)

type jobs struct{ rows []invitation.Delivery; err error }
func (j jobs) ListDeliveryPage(_ context.Context, tenant, after string, limit int)([]invitation.Delivery,error){
	if tenant!="school"{return nil,errors.New("incorrect tenant")}
	if j.err!=nil{return nil,j.err}
	var out []invitation.Delivery
	for _,row:=range j.rows{if row.DeliveryID>after&&len(out)<limit{out=append(out,row)}}
	return out,nil
}
type receipts []maildelivery.Receipt
func (r receipts) ListReceiptPage(_ context.Context,tenant,after string,limit int)([]maildelivery.Receipt,error){
	if tenant!="school"{return nil,errors.New("incorrect tenant")}
	var out []maildelivery.Receipt
	for _,row:=range r{if row.DeliveryID>after&&len(out)<limit{out=append(out,row)}}
	return out,nil
}
func id(n int)string{return fmt.Sprintf("00000000-0000-0000-0000-%012d",n)}
func TestReconciliationPages(t *testing.T){
	now:=time.Date(2026,11,28,0,0,0,0,time.UTC)
	old:=now.Add(-3*time.Minute)
	var source jobs
	var receipt receipts
	for n:=1;n<=260;n++{
		source.rows=append(source.rows,invitation.Delivery{DeliveryID:id(n),AccountID:"account",Purpose:"ACCOUNT_SETUP",SentAt:&now})
		receipt=append(receipt,maildelivery.Receipt{DeliveryID:id(n),State:"ACCEPTED",AcceptedAt:&now})
	}
	source.rows=append(source.rows,invitation.Delivery{DeliveryID:id(262),AccountID:"account",Purpose:"PASSWORD_RESET"},invitation.Delivery{DeliveryID:id(263),AccountID:"account",Purpose:"ACCOUNT_SETUP"})
	receipt=append(receipt,maildelivery.Receipt{DeliveryID:id(261),State:"ACCEPTED"},maildelivery.Receipt{DeliveryID:id(262),State:"ACCEPTED"},maildelivery.Receipt{DeliveryID:id(263),State:"DISPATCHING",ReservedAt:&old})
	snapshot:=func(ctx context.Context,read func(context.Context,time.Time)error)error{return read(ctx,now)}
	op:=impl.New(source,receipt,snapshot)
	first,err:=op.Search(context.Background(),mail.SearchRequest{TenantID:"school",Limit:2})
	if err!=nil||len(first.Rows)!=2||first.NextCursor!=id(262){t.Fatalf("first page: %+v %v",first,err)}
	if first.Rows[0].Disposition!="RECEIPT_WITHOUT_JOB"||first.Rows[0].AccountID!=nil||first.Rows[1].Disposition!="ACCEPTED_WITHOUT_SENT"{t.Fatal(first)}
	second,err:=op.Search(context.Background(),mail.SearchRequest{TenantID:"school",After:first.NextCursor,Limit:2})
	if err!=nil||len(second.Rows)!=1||second.NextCursor!=""||second.Rows[0].Disposition!="STALE_DISPATCH"{t.Fatal(second,err)}
	for _,input:=range []mail.SearchRequest{{TenantID:"",Limit:1},{TenantID:"school",Limit:0},{TenantID:"school",Limit:501},{TenantID:"school",Limit:1,After:"invalid"}}{
		if _,err:=op.Search(context.Background(),input);!errors.Is(err,mail.ErrInvalidScope){t.Fatal(input,err)}
	}
	failed:=impl.New(jobs{err:errors.New("unavailable")},receipt,snapshot)
	if report,err:=failed.Search(context.Background(),mail.SearchRequest{TenantID:"school",Limit:1});err==nil||report.TenantID!=""{t.Fatal(report,err)}
}
