package operator

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"

type SchoolCheckParams struct {
	SchoolID string `path:"schoolId"`
}

type SchoolCheckRequest struct {
	request.Request
	Params SchoolCheckParams
}
