export interface OperatorActivity {
  correlationId?: string;
  eventId: string;
  actorId: string;
  resourceId: string;
  operation: string;
  occurredAt: string;
}
export interface DeliverySummary {
  enabled: boolean;
  health: 'DISABLED' | 'HEALTHY' | 'DEGRADED';
  rejectedLast24Hours: number;
  limitedLast24Hours: number;
  pending: number;
  blocked: number;
  delivered: number;
  oldestPendingAt: string | null;
  lastDeliveredAt: string | null;
}
export interface BlockedEvent {
  eventId: string;
  actorId: string;
  resourceId: string;
  action: string;
  attempts: number;
  blockedAt: string;
}

export interface SchoolReceipt {
  operation: 'RENAME' | 'ACTIVATE' | 'DISABLE';
  executionId: string;
  schoolId: string;
  previousVersion: string;
  correlationId?: string;
  occurredAt: string;
}
