CREATE TABLE internal_tool.environments (
 id uuid PRIMARY KEY, code text UNIQUE NOT NULL, name text NOT NULL,
 created_by uuid NOT NULL REFERENCES internal_tool.operators(id),
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE internal_tool.api_templates (
 id uuid PRIMARY KEY, code text NOT NULL, version integer NOT NULL CHECK(version>0),
 module text NOT NULL, method text NOT NULL CHECK(method IN ('GET','POST','PATCH','DELETE')),
 endpoint text NOT NULL, scope text NOT NULL CHECK(scope IN ('GLOBAL','SCHOOL')),
 headers jsonb NOT NULL, body_schema jsonb NOT NULL,
 created_by uuid NOT NULL REFERENCES internal_tool.operators(id),
 created_at timestamptz NOT NULL DEFAULT now(), UNIQUE(code,version)
);
CREATE TABLE internal_tool.approval_requests (
 id uuid PRIMARY KEY, template_id uuid NOT NULL REFERENCES internal_tool.api_templates(id),
 environment_id uuid NOT NULL REFERENCES internal_tool.environments(id),
 school_id uuid, input jsonb NOT NULL, reason text NOT NULL,
 requester_id uuid NOT NULL REFERENCES internal_tool.operators(id),
 status text NOT NULL CHECK(status IN ('PENDING','APPROVED','REJECTED','CANCELLED')),
 version bigint NOT NULL DEFAULT 1 CHECK(version>0),
 reviewer_id uuid REFERENCES internal_tool.operators(id), decision_comment text,
 created_at timestamptz NOT NULL DEFAULT now(), decided_at timestamptz,
 CHECK ((status='PENDING' AND reviewer_id IS NULL AND decided_at IS NULL) OR
 (status<>'PENDING' AND reviewer_id IS NOT NULL AND decided_at IS NOT NULL))
);
CREATE INDEX internal_tool_approval_inbox ON internal_tool.approval_requests(status,created_at DESC,id);
CREATE INDEX internal_tool_approval_owner ON internal_tool.approval_requests(requester_id,created_at DESC,id);
CREATE TABLE internal_tool.workflow_events (
 id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 actor_id uuid NOT NULL REFERENCES internal_tool.operators(id),
 resource_id uuid NOT NULL, action text NOT NULL, occurred_at timestamptz NOT NULL DEFAULT now()
);
