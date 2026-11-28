export const SCHOOL_MUTATION_OPERATIONS = ['RENAME', 'ACTIVATE', 'DISABLE'] as const;
export function validSchoolId(value: string): boolean {
 return value === '' || (/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/.test(value) && value !== '00000000-0000-0000-0000-000000000000');
}
export function validOperation(value: string): boolean {
 return value === '' || SCHOOL_MUTATION_OPERATIONS.some((operation) => operation === value);
}
export function buildReceiptFilters(schoolId: string, operation: string): Record<string, string> | null {
 if (!validSchoolId(schoolId) || !validOperation(operation)) return null;
 return { ...(schoolId && { schoolId }), ...(operation && { operation }) };
}
