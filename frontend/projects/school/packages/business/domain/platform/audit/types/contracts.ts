export interface AuditActivity {
  auditId: string;
  source: string;
  assetType: string;
  assetId: string;
  operation: string;
  occurredAt: string;
  actorId?: string;
  data: unknown;
}
