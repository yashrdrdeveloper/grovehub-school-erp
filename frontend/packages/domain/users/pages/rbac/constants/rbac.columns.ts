import { Column } from '@Components/components/Table';
import type { Role } from '@Business/domain/users/rbac/types/contracts';

export const RBAC_COLUMNS = Object.freeze([
  new Column<Role>().setId('name').setHeader(__('Name')).setAccessor('name').setWidth(180).setSortable(true).build(),
  new Column<Role>().setId('status').setHeader(__('Status')).setAccessor('status').setWidth(180).setSortable(true).build(),
  new Column<Role>().setId('description').setHeader(__('Description')).setAccessor('description').setWidth(180).build(),
]);
