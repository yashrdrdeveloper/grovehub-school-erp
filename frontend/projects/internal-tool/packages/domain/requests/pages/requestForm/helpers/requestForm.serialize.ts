import type { Field } from '../../../../../business/api/types';

function invalid(name: string): never {
  throw new Error('Invalid schema field: ' + name);
}

function serializeValue(name: string, field: Field, raw: FormDataEntryValue | null): string {
  if (typeof raw !== 'string') return invalid(name);
  if (field.type === 'string') return JSON.stringify(raw);
  if (field.type === 'boolean') {
    if (raw === 'true' || raw === 'false') return raw;
    return invalid(name);
  }
  const text = raw.trim();
  if (field.type === 'integer') {
    if (!/^-?(0|[1-9][0-9]*)$/.test(text)) return invalid(name);
    const value = BigInt(text);
    if (value < -9223372036854775808n || value > 9223372036854775807n) return invalid(name);
    return text;
  }
  if (field.type === 'number') {
    if (text.length > 128 || !/^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$/.test(text)) return invalid(name);
    return text;
  }
  const value: unknown = JSON.parse(text);
  if (field.type === 'array' && Array.isArray(value)) return text;
  if (field.type === 'object' && value !== null && typeof value === 'object' && !Array.isArray(value)) return text;
  return invalid(name);
}

export function serializeSchemaInput(schema: Record<string, Field>, data: FormData): string {
  const entries = Object.entries(schema).flatMap(([name, field]) => {
    const mode = data.get('mode:' + name);
    if (mode === 'omit' && !field.required) return [];
    if (mode === 'null' && field.nullable) return [JSON.stringify(name) + ':null'];
    if (mode !== 'value') return invalid(name);
    return [JSON.stringify(name) + ':' + serializeValue(name, field, data.get('value:' + name))];
  });
  return '{' + entries.join(',') + '}';
}

export function serializeRequest(templateId: string, schoolScoped: boolean, schema: Record<string, Field>, data: FormData): string {
  const schoolId = schoolScoped && data.get('schoolId');
  const metadata = JSON.stringify({
    templateId,
    environmentId: data.get('environmentId'),
    schoolId: schoolId || null,
    reason: data.get('reason'),
  });
  return metadata.slice(0, -1) + ',"input":' + serializeSchemaInput(schema, data) + '}';
}
