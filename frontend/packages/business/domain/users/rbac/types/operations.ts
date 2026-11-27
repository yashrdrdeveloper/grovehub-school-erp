export interface RoleMetadata { name: string; description: string; personas: string[] }
export interface PermissionIntent { roleId: string; version: number; permissionCodes: string[] }
export interface PermissionPreview { roleId: string; version: number; added: string[]; removed: string[]; approvalRequired: boolean }
export interface Assignment { accountId: string; version: number; primaryRoleId: string; secondaryRoleIds: string[] }
export interface AssignmentPreview extends Assignment { roleVersions: { roleId: string; version: number }[]; approvalRequired: boolean }
export interface Capabilities { permissions: string[]; actions: Record<string, boolean>; executionConfigured: boolean; executionRequiresReauthentication: boolean }
export interface PermissionCatalog { permissions: { code: string; name: string; feature: string; canDelegate: boolean }[] }
export interface Approval { requestId: string; requestType: string; status: string; version: number; requesterId: string; createdAt: string }
export interface ApprovalDetail extends Approval { intent: unknown; reason: string }
