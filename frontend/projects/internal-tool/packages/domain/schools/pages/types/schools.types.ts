import type { Environment, SchoolGrant, Template } from '../../../../business/api/types';
export interface SchoolAction {
 readonly kind: 'rename' | 'lifecycle';
 readonly template: Template;
 readonly grant: SchoolGrant;
}
export interface SchoolRecord {
 readonly key: string;
 readonly schoolId: string;
 readonly environment: Environment;
 readonly name: string;
 readonly code: string;
 readonly version: string;
 readonly status?: 'ACTIVE' | 'DISABLED';
 readonly snapshotReady: boolean;
 readonly actions: readonly SchoolAction[];
}
