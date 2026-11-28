package helpers_test

import (
	"testing"
	"time"

	"grovehub/services/identity/client/invitation"
	"grovehub/services/integration/internal/app/repository/maildelivery"
	"grovehub/services/integration/internal/app/services/mail/operations/reconciliation/helpers"
)
func TestClassify(t *testing.T){
	now:=time.Now()
	old:=now.Add(-3*time.Minute)
	reason:="DELIVERY_UNCERTAIN"
	for _,tc:=range []struct{name string;job *invitation.Delivery;receipt *maildelivery.Receipt;want string}{
		{"uncertain",&invitation.Delivery{SuppressionReason:&reason},nil,"REVIEW_UNCERTAIN"},
		{"sent without receipt acceptance",&invitation.Delivery{SentAt:&now},&maildelivery.Receipt{State:"DISPATCHING"},"SENT_WITHOUT_ACCEPTANCE"},
		{"stale",&invitation.Delivery{},&maildelivery.Receipt{State:"DISPATCHING",ReservedAt:&old},"STALE_DISPATCH"},
		{"active lease",&invitation.Delivery{LeaseUntil:&now},&maildelivery.Receipt{State:"DISPATCHING",ReservedAt:&old},""},
		{"consistent",&invitation.Delivery{SentAt:&now},&maildelivery.Receipt{State:"ACCEPTED"},""},
		{"receipt only",nil,&maildelivery.Receipt{State:"ACCEPTED"},"RECEIPT_WITHOUT_JOB"},
		{"pending",&invitation.Delivery{},nil,""},
	}{
		t.Run(tc.name,func(t *testing.T){if row:=helpers.Classify(helpers.Pair{ID:"id",Job:tc.job,Receipt:tc.receipt},now);row.Disposition!=tc.want{t.Fatal(row)}})
	}
}
