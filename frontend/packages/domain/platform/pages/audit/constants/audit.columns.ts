import { Column } from '@Components/components/Table';
import type { AuditActivity } from '@Business/domain/platform/audit/types/contracts';

export const AUDIT_COLUMNS = Object.freeze([
  new Column<AuditActivity>().setId('operation').setHeader(__('Operation')).setAccessor('operation').setWidth(180).setSortable(true).build(),
  new Column<AuditActivity>().setId('source').setHeader(__('Source')).setAccessor('source').setWidth(180).setSortable(true).build(),
  new Column<AuditActivity>().setId('assetType').setHeader(__('Asset type')).setAccessor('assetType').setWidth(180).setSortable(true).build(),
  new Column<AuditActivity>().setId('occurredAt').setHeader(__('Occurred at')).setAccessor('occurredAt').setWidth(180).setSortable(true).build(),
]);
