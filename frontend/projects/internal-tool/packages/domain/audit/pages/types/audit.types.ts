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
