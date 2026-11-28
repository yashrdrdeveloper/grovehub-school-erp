// Package activity records local workflow occurrences in the caller's transaction.
package activity

import (
 "context"
 "database/sql"

 corehttp "grovehub/libs/core/http"
)

type Executor interface{
 ExecContext(context.Context,string,...any)(sql.Result,error)
}

func Record(ctx context.Context,db Executor,actor,resource,operation string)error{
 ctx=corehttp.Ensure(ctx)
 _,err:=db.ExecContext(ctx,"INSERT INTO internal_tool.workflow_events(actor_id,resource_id,action,correlation_id) VALUES($1,$2,$3,$4)",actor,resource,operation,corehttp.ID(ctx))
 return err
}
