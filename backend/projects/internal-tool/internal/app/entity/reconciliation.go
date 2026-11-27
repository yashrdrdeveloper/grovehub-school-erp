package entity

import "time"

type Reconciliation struct{
 ExecutionID string `json:"executionId"`
 ReviewerID string `json:"reviewerId"`
 Outcome string `json:"outcome"`
 Evidence string `json:"evidence"`
 CreatedAt time.Time `json:"createdAt"`
}
