import { FilterFieldType, FilterOperator } from '@Components/components/Filters';
import type { FilterField } from '@Components/components/Filters';

export const AUDIT_FILTERS: readonly FilterField[] = Object.freeze([
  {
    id: 'operation',
    label: __('Operation'),
    type: FilterFieldType.TEXT,
    operators: [FilterOperator.EQUALS],
  },
  {
    id: 'source',
    label: __('Source'),
    type: FilterFieldType.TEXT,
    operators: [FilterOperator.EQUALS],
  },
  {
    id: 'assetType',
    label: __('Asset type'),
    type: FilterFieldType.TEXT,
    operators: [FilterOperator.EQUALS],
  },
]);
