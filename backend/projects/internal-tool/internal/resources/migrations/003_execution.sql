CREATE TABLE internal_tool.executions (
 id uuid PRIMARY KEY,
 request_id uuid NOT NULL UNIQUE REFERENCES internal_tool.approval_requests(id),
 operator_id uuid NOT NULL REFERENCES internal_tool.operators(id),
 status text NOT NULL CHECK(status IN ('RUNNING','SUCCEEDED','FAILED','UNKNOWN')),
 response_status integer CHECK(response_status BETWEEN 200 AND 599),
 created_at timestamptz NOT NULL DEFAULT now(),
 finished_at timestamptz,
 CHECK ((status='RUNNING' AND finished_at IS NULL) OR (status<>'RUNNING' AND finished_at IS NOT NULL))
);
