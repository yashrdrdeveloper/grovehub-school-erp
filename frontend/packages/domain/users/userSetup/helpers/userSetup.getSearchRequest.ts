import { FilterOperator } from '@Components/components/Filters';
import { _map } from '@Base/utilities/map';
import type { SearchRequest } from '@Business/domain/users';

import { UserField } from '../constants/userSetup.fields';
import { getUserSetupDateBoundary } from './userSetup.getDateBoundary';
import type { UserSetupQueryState } from '../types/userSetup.state';

export const getUserSetupSearchRequest = (state: UserSetupQueryState): SearchRequest => {
  const quickStatus = state.quickFilterIds[0];
  const filters: NonNullable<SearchRequest['filters']> = state.filters.flatMap((item) => {
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
      return [{ field: item.field, operator: operators[item.operator as keyof typeof operators], values: [...item.values] }];
    }
    return [];
  });
  if (quickStatus && quickStatus !== 'ALL') {
    filters.push({ field: UserField.STATUS, operator: 'EQUALS', values: [quickStatus] });
  }
  return {
    searchText: state.searchText.trim(),
    filters,
    sort: _map(state.sorting, (rule) => ({ field: rule.field, direction: rule.direction })),
    pageInfo: { rows: state.pageSize, cursor: state.cursors[state.cursors.length - 1] },
  };
};
