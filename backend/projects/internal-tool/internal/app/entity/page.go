package entity

import "time"

// PageRequest supports bounded page-number navigation. RowSize defaults to 25.
type PageRequest struct{
 PageNumber int `json:"pageNumber"`
 RowSize int `json:"rowSize"`
}
func(p PageRequest) Normalize()(PageRequest,bool){
 if p.PageNumber==0{p.PageNumber=1}
 if p.RowSize==0{p.RowSize=25}
 return p,p.PageNumber>=1&&p.PageNumber<=10000&&p.RowSize>=1&&p.RowSize<=100
}
func(p PageRequest) Offset()int{return (p.PageNumber-1)*p.RowSize}
type Page[T any] struct{
 Records []T `json:"records"`
 PageNumber int `json:"pageNumber"`
 RowSize int `json:"rowSize"`
 HasMore bool `json:"hasMore"`
}
func NewPage[T any](p PageRequest,records []T)Page[T]{
 more:=len(records)>p.RowSize
 if more{records=records[:p.RowSize]}
 if records==nil{records=[]T{}}
 return Page[T]{Records:records,PageNumber:p.PageNumber,RowSize:p.RowSize,HasMore:more}
}
type WorkflowEvent struct{
 ID int64 `json:"id"`
 ActorID string `json:"actorId"`
 ResourceID string `json:"resourceId"`
 Action string `json:"action"`
 OccurredAt time.Time `json:"occurredAt"`
}
