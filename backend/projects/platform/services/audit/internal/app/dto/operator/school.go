package operator

import "grovehub/libs/core/request"

type SchoolCheckParams struct {
	SchoolID string `path:"schoolId"`
}

type SchoolCheckRequest struct {
	request.Request
	Params SchoolCheckParams
}
