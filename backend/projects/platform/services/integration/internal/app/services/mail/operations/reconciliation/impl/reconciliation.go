package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"grovehub/services/integration/client/mail"
	"grovehub/services/integration/internal/app/services/mail/operations/reconciliation"
	"grovehub/services/integration/internal/app/services/mail/operations/reconciliation/helpers"
)

type operation struct {
	jobs reconciliation.JobReader
	receipts reconciliation.ReceiptReader
	snapshot reconciliation.Snapshot
}

func New(jobs reconciliation.JobReader, receipts reconciliation.ReceiptReader, snapshot reconciliation.Snapshot) reconciliation.Operation {
	return &operation{jobs:jobs,receipts:receipts,snapshot:snapshot}
}

func (o *operation) Search(ctx context.Context, in mail.SearchRequest) (mail.Report,error) {
	if in.TenantID=="" || strings.TrimSpace(in.TenantID)!=in.TenantID || in.Limit<1 || in.Limit>500 {
		return mail.Report{},mail.ErrInvalidScope
	}
	if in.After=="" {in.After=uuid.Nil.String()}
	id,err:=uuid.Parse(in.After)
	if err!=nil||id.String()!=in.After{return mail.Report{},mail.ErrInvalidScope}
	if o.jobs==nil||o.receipts==nil||o.snapshot==nil{return mail.Report{},errors.New("reconciliation unavailable")}
	report:=mail.Report{TenantID:in.TenantID,Rows:make([]mail.Row,0,in.Limit+1)}
	err=o.snapshot(ctx,func(ctx context.Context, now time.Time)error{
		after:=in.After
		for {
			if err:=ctx.Err();err!=nil{return err}
			jobs,err:=o.jobs.ListDeliveryPage(ctx,in.TenantID,after,reconciliation.BatchSize);if err!=nil{return err}
			receipts,err:=o.receipts.ListReceiptPage(ctx,in.TenantID,after,reconciliation.BatchSize);if err!=nil{return err}
			pairs:=helpers.Merge(jobs,receipts,reconciliation.BatchSize)
			for _,pair:=range pairs{
				if pair.ID<=after{return errors.New("invalid reconciliation page order")}
				after=pair.ID
				row:=helpers.Classify(pair,now)
				if row.Disposition==""{continue}
				report.Rows=append(report.Rows,row)
				if len(report.Rows)>in.Limit{return nil}
			}
			if len(pairs)<reconciliation.BatchSize{return nil}
		}
	})
	if err!=nil{return mail.Report{},err}
	if len(report.Rows)>in.Limit{
		report.Rows=report.Rows[:in.Limit]
		report.NextCursor=report.Rows[in.Limit-1].DeliveryID
	}
	return report,nil
}
