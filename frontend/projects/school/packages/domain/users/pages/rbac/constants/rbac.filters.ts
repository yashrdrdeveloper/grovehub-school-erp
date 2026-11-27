import { FilterFieldType, FilterOperator } from '@Components/components/Filters';
import type { FilterField } from '@Components/components/Filters';

export const RBAC_FILTERS: readonly FilterField[] = Object.freeze([
  { id: 'name', label: __('Name'), type: FilterFieldType.TEXT, operators: [FilterOperator.EQUALS] },
  { id: 'status', label: __('Status'), type: FilterFieldType.TEXT, operators: [FilterOperator.EQUALS] },
]);
