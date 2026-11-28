import { apiRequest } from '@Business/api';
import { HTTP_METHOD } from '@Business/api/constants/api.method';
import type { SearchRequest, SearchResult } from '@Business/accounts/types';
import type { Role } from '../types/contracts';
import type {
  RoleMetadata,
  PermissionIntent,
  PermissionPreview,
  Assignment,
  AssignmentPreview,
  Capabilities,
  PermissionCatalog,
  Approval,
  ApprovalDetail,
} from '../types/operations';

const root = '/identity/protected';
const rolePath = (id: string) => root + '/roles/' + encodeURIComponent(id);
const assignmentPath = (id: string) => root + '/account-roles/' + encodeURIComponent(id);
const approvalPath = (id: string) => root + '/approvals/' + encodeURIComponent(id);
const submit = <T>(path: string, body: unknown, key: string) =>
  apiRequest<T>(path, { method: HTTP_METHOD.POST, body, headers: { 'Idempotency-Key': key } });
export const rbacOperations = {
  capabilities: () => apiRequest<Capabilities>(root + '/capabilities/me'),
  catalog: () => apiRequest<PermissionCatalog>(root + '/capabilities/permissions'),
  create: (body: RoleMetadata) =>
    apiRequest<Role>(root + '/roles', { method: HTTP_METHOD.POST, body }),
  update: (id: string, body: RoleMetadata & { version: number }) =>
    apiRequest<Role>(rolePath(id), { method: HTTP_METHOD.PATCH, body }),
  previewPermissions: ({ roleId, ...body }: PermissionIntent) =>
    apiRequest<PermissionPreview>(rolePath(roleId) + '/permissions/preview', {
      method: HTTP_METHOD.POST,
      body,
    }),
  submitPermissions: (intent: PermissionIntent, reason: string, key: string) =>
    submit<Approval>(root + '/approvals/permissions', { intent, reason }, key),
  assignment: (id: string) => apiRequest<Assignment>(assignmentPath(id)),
  previewAssignment: ({ accountId, ...body }: Assignment) =>
    apiRequest<AssignmentPreview>(assignmentPath(accountId) + '/preview', {
      method: HTTP_METHOD.POST,
      body,
    }),
  submitAssignment: (preview: AssignmentPreview, reason: string, key: string) => {
    const { accountId, version, primaryRoleId, secondaryRoleIds, roleVersions } = preview;
    return submit<Approval>(
      root + '/approvals/assignments',
      { intent: { accountId, version, primaryRoleId, secondaryRoleIds, roleVersions }, reason },
      key,
    );
  },
  history: (body: SearchRequest) =>
    apiRequest<SearchResult<Approval>>(root + '/approvals/history/search', {
      method: HTTP_METHOD.POST,
      body,
    }),
  inbox: (body: SearchRequest) =>
    apiRequest<SearchResult<Approval>>(root + '/approvals/inbox/search', {
      method: HTTP_METHOD.POST,
      body,
    }),
  detail: (id: string) => apiRequest<ApprovalDetail>(approvalPath(id)),
  review: (id: string, action: 'approve' | 'reject', version: number, comment: string) =>
    apiRequest<Approval>(approvalPath(id) + '/' + action, {
      method: HTTP_METHOD.POST,
      body: { version, comment },
    }),
  cancel: (id: string, version: number) =>
    apiRequest<Approval>(approvalPath(id) + '/cancel', {
      method: HTTP_METHOD.POST,
      body: { version },
    }),
  execute: (id: string, version: number) =>
    apiRequest<Approval>(approvalPath(id) + '/execute', {
      method: HTTP_METHOD.POST,
      body: { version },
    }),
};
