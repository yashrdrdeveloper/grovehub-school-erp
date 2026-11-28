import type { Environment, SchoolGrant, Template } from '../../../../business/api/types';
import type { SchoolAction, SchoolRecord } from '../types/schools.types';

function actionKind(template: Template, grant: SchoolGrant): SchoolAction['kind'] | undefined {
 if (!grant.canRequest || template.scope !== 'SCHOOL' || template.method !== 'PATCH') return undefined;
 if (template.endpoint === '/school/restricted/registry/schools/${schoolId}/name' && grant.schoolStatus === 'ACTIVE') return 'rename';
 if (template.endpoint === '/school/restricted/registry/schools/${schoolId}/status' && grant.lifecycle) return 'lifecycle';
 return undefined;
}

export function buildSchoolRecords(grants: readonly SchoolGrant[], environments: readonly Environment[], templates: readonly Template[]): SchoolRecord[] {
 const groups = new Map<string, SchoolGrant[]>();
 for (const grant of grants) {
  if (!grant.canRequest && !grant.canReview && !grant.canExecute) continue;
  const key = JSON.stringify([grant.environmentCode, grant.schoolId]);
  const group = groups.get(key) ?? [];
  groups.set(key, [...group, grant]);
 }
 return [...groups.entries()].flatMap(([key, group]) => {
  const first = group[0];
  if (!first) return [];
  const environment = environments.find((item) => item.code === first.environmentCode);
  if (!environment) return [];
  const snapshotReady = Boolean(first.schoolName && first.schoolCode && first.schoolVersion && /^[1-9][0-9]*$/.test(first.schoolVersion) && (first.schoolStatus === 'ACTIVE' || first.schoolStatus === 'DISABLED') && group.every((grant) => grant.schoolName === first.schoolName && grant.schoolCode === first.schoolCode && grant.schoolVersion === first.schoolVersion && grant.schoolStatus === first.schoolStatus));
  const actions: SchoolAction[] = [];
  if (snapshotReady) {
   for (const grant of group) {
    const template = templates.find((item) => item.id === grant.templateId);
    if (!template) continue;
    const kind = actionKind(template, grant);
    if (!kind || actions.some((item) => item.template.id === template.id)) continue;
    actions.push({ kind, template, grant });
   }
  }
  return [{ key, schoolId: first.schoolId, environment, name: first.schoolName ?? first.schoolId, code: first.schoolCode ?? '', version: first.schoolVersion ?? '', status: first.schoolStatus, snapshotReady, actions }];
 }).sort((a,b) => a.name.localeCompare(b.name) || a.environment.name.localeCompare(b.environment.name) || a.key.localeCompare(b.key));
}

export function filterSchools(records: readonly SchoolRecord[], query: string, environmentId: string, status: string): SchoolRecord[] {
 const normalized = query.trim().toLowerCase();
 return records.filter((record) => (!environmentId || record.environment.id === environmentId) &&
  (!status || (record.snapshotReady && record.status === status)) &&
  (!normalized || [record.name, record.code, record.schoolId].some((value) => value.toLowerCase().includes(normalized))));
}
