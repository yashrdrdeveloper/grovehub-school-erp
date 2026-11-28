import type { Environment, SchoolGrant, Template } from '../../../../../business/api/types';

export function getRequestTemplates(templates: readonly Template[], grants: readonly SchoolGrant[]) {
  return templates.filter((template) => template.scope === 'GLOBAL' ||
    grants.some((grant) => grant.canRequest && grant.templateId === template.id));
}

export function getRequestEnvironments(template: Template | undefined, environments: readonly Environment[], grants: readonly SchoolGrant[]) {
  if (!template) return [];
  if (template.scope === 'GLOBAL') return environments;
  return environments.filter((environment) => grants.some((grant) =>
    grant.canRequest && grant.templateId === template.id && grant.environmentCode === environment.code));
}

export function getRequestSchoolIds(template: Template | undefined, environment: Environment | undefined, grants: readonly SchoolGrant[]) {
  if (!template || template.scope !== 'SCHOOL' || !environment) return [];
  return [...new Set(grants.filter((grant) => grant.canRequest &&
    grant.templateId === template.id && grant.environmentCode === environment.code).map((grant) => grant.schoolId))].sort();
}
