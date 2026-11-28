package entity

import (
 "encoding/json"
 "time"
)

type Environment struct {
 ID string `json:"id"`
 Code string `json:"code"`
 Name string `json:"name"`
}
type Template struct {
 ID string `json:"id"`
 Code string `json:"code"`
 Version int `json:"version"`
 Module string `json:"module"`
 Method string `json:"method"`
 Endpoint string `json:"endpoint"`
 Scope string `json:"scope"`
 Headers map[string]string `json:"headers"`
 BodySchema map[string]Field `json:"bodySchema"`
}
type Field struct {
 Type string `json:"type"`
 Required bool `json:"required"`
 Properties map[string]Field `json:"properties,omitempty"`
 Items *Field `json:"items,omitempty"`
 Nullable bool `json:"nullable,omitempty"`
}
type Approval struct {
 ID string `json:"id"`
 TemplateID string `json:"templateId"`
 EnvironmentID string `json:"environmentId"`
 SchoolID *string `json:"schoolId"`
 Input json.RawMessage `json:"input"`
 Reason string `json:"reason"`
 RequesterID string `json:"requesterId"`
 Status string `json:"status"`
 Version int64 `json:"version"`
 ReviewerID *string `json:"reviewerId"`
 DecisionComment *string `json:"decisionComment"`
 CreatedAt time.Time `json:"createdAt"`
 DecidedAt *time.Time `json:"decidedAt"`
}
