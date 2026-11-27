package entity

import "time"

type Execution struct {
 ID string `json:"id"`
 RequestID string `json:"requestId"`
 OperatorID string `json:"operatorId"`
 Status string `json:"status"`
 ResponseStatus *int `json:"responseStatus"`
 CreatedAt time.Time `json:"createdAt"`
 FinishedAt *time.Time `json:"finishedAt"`
}
