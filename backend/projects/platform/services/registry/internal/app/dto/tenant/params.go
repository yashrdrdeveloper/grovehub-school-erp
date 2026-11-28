package tenant

type Params struct {
	ProjectID string `path:"projectId"`
	TenantID string `path:"tenantId"`
}

type ProjectParams struct {
	ProjectID string `path:"projectId"`
}
