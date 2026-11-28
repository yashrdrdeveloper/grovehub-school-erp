package helpers

import (
	"time"

	"grovehub/services/identity/client/invitation"
	"grovehub/services/integration/client/mail"
	"grovehub/services/integration/internal/app/repository/maildelivery"
)

type Pair struct {
	ID string
	Job *invitation.Delivery
	Receipt *maildelivery.Receipt
}

// Merge returns the first limit distinct IDs from two sorted pages. Neither
// source can hide an earlier ID beyond a page of at least limit rows.
func Merge(jobs []invitation.Delivery, receipts []maildelivery.Receipt, limit int) []Pair {
	pairs := make([]Pair,0)
	i,j:=0,0
	for len(pairs)<limit && (i<len(jobs)||j<len(receipts)) {
		var pair Pair
		if j==len(receipts)||(i<len(jobs)&&jobs[i].DeliveryID<receipts[j].DeliveryID) {
			pair=Pair{ID:jobs[i].DeliveryID,Job:&jobs[i]};i++
		} else if i==len(jobs)||receipts[j].DeliveryID<jobs[i].DeliveryID {
			pair=Pair{ID:receipts[j].DeliveryID,Receipt:&receipts[j]};j++
		} else {
			pair=Pair{ID:jobs[i].DeliveryID,Job:&jobs[i],Receipt:&receipts[j]};i++;j++
		}
		pairs=append(pairs,pair)
	}
	return pairs
}

func Classify(pair Pair, now time.Time) mail.Row {
	row:=mail.Row{DeliveryID:pair.ID,IdentityState:"MISSING",IntegrationState:"MISSING"}
	job,receipt:=pair.Job,pair.Receipt
	if receipt!=nil {
		row.IntegrationState=receipt.State
		row.ReservedAt=receipt.ReservedAt
		row.AcceptedAt=receipt.AcceptedAt
	}
	if job==nil {row.Disposition="RECEIPT_WITHOUT_JOB";return row}
	row.AccountID=&job.AccountID
	row.Purpose=&job.Purpose
	row.SuppressionReason=job.SuppressionReason
	row.IdentityState="PENDING"
	if job.SentAt!=nil {row.IdentityState="SENT"} else if job.SuppressedAt!=nil {row.IdentityState="SUPPRESSED"}
	switch {
	case job.SuppressionReason!=nil && *job.SuppressionReason=="DELIVERY_UNCERTAIN":
		row.Disposition="REVIEW_UNCERTAIN"
	case receipt!=nil && receipt.State=="ACCEPTED" && job.SentAt==nil:
		row.Disposition="ACCEPTED_WITHOUT_SENT"
	case receipt!=nil && receipt.State=="DISPATCHING" && job.SentAt!=nil:
		row.Disposition="SENT_WITHOUT_ACCEPTANCE"
	case receipt!=nil && receipt.State=="DISPATCHING" && receipt.ReservedAt!=nil &&
		receipt.ReservedAt.Before(now.Add(-2*time.Minute)) && (job.LeaseUntil==nil||job.LeaseUntil.Before(now)):
		row.Disposition="STALE_DISPATCH"
	}
	return row
}
