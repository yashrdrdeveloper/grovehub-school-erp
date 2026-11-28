package workload

type Params struct {
	WorkloadID string `path:"workloadId"`
}

type KeyParams struct {
	WorkloadID string `path:"workloadId"`
	KeyID      string `path:"keyId"`
}

type GrantParams struct {
	WorkloadID string `path:"workloadId"`
	GrantID    string `path:"grantId"`
}
