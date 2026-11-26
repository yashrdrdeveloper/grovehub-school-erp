import { FilterOperator } from '@Components/components/Filters';
import { _map } from '@Base/utilities/map';
import { _flatMap } from '@Base/utilities/flatMap';
import { _trim } from '@Base/utilities/trim';
import type { SearchRequest } from '@Business/domain/users';

import { UserField } from '../constants/userSetup.fields';
import { getUserSetupDateBoundary } from './userSetup.getDateBoundary';
import type { UserSetupQueryState } from '../types/userSetup.state';

export const getUserSetupSearchRequest = (state: UserSetupQueryState): SearchRequest => {
  const quickStatus = state.quickFilterIds[0];
  const filters = _flatMap(state.filters, (item): NonNullable<SearchRequest['filters']> => {
    if (item.operator === FilterOperator.BETWEEN) {
      const start = getUserSetupDateBoundary(item.values[0], 0);
      const end = getUserSetupDateBoundary(item.values[1], 1);
      if (!start || !end || start >= end) return [];
      return [
        { field: item.field, operator: 'GTE', values: [start] },
        { field: item.field, operator: 'LT', values: [end] },
      ];
    }
    const operators = {
      [FilterOperator.EQUALS]: 'EQUALS',
      [FilterOperator.NOT_EQUALS]: 'NOT_EQUALS',
      [FilterOperator.IN]: 'IN',
      [FilterOperator.NOT_IN]: 'NIN',
    } as const;
    if (item.operator in operators) {
      return [
        {
          field: item.field,
          operator: operators[item.operator as keyof typeof operators],
          values: [...item.values],
        },
      ];
    }
    return [];
  });
  if (quickStatus && quickStatus !== 'ALL') {
    filters.push({ field: UserField.STATUS, operator: 'EQUALS', values: [quickStatus] });
  }
  return {
    searchText: _trim(state.searchText),
    filters,
    sort: _map(state.sorting, (rule) => ({ field: rule.field, direction: rule.direction })),
    pageInfo: {
      pageNumber: state.pageNumber,
      rowSize: state.pageSize,
      cursor: state.cursors[state.pageNumber],
    },
  };
};
