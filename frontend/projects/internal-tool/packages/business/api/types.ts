export type Persona = 'L1' | 'L2';
export interface Operator { id: string; email: string; persona: Persona; status?: string }
export interface Environment { id: string; code: string; name: string }
export interface Field { type: string; required: boolean; nullable?: boolean; properties?: Record<string, Field>; items?: Field }
export interface Template { id: string; code: string; version: number; module: string; method: string; endpoint: string; scope: 'GLOBAL' | 'SCHOOL'; headers: Record<string, string>; bodySchema: Record<string, Field> }
export interface Approval { id: string; templateId: string; environmentId: string; schoolId: string | null; input: unknown; reason: string; requesterId: string; status: string; version: number; reviewerId: string | null; decisionComment: string | null; createdAt: string; decidedAt: string | null }
export interface Execution { id: string; requestId: string; operatorId: string; status: string; responseStatus: number | null; createdAt: string; finishedAt: string | null }
export interface Reconciliation { executionId: string; reviewerId: string; outcome: string; evidence: string; createdAt: string }
export interface WorkflowEvent { correlationId?: string; id: number; actorId: string; resourceId: string; action: string; occurredAt: string }
export interface Page<T> { records: T[]; pageNumber: number; rowSize: number; hasMore: boolean }

export interface SchoolGrant {
  environmentCode: string;
  templateId: string;
  schoolId: string;
  canRequest: boolean;
  canReview: boolean;
  canExecute: boolean;
}
