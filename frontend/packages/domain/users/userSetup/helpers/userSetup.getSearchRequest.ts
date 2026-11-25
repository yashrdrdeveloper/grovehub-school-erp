import { _map } from '@Base/utilities/map';
import type { SearchRequest } from '@Business/domain/users';

import { UserField } from '../constants/userSetup.fields';
import type { UserSetupQueryState } from '../types/userSetup.state';

export const getUserSetupSearchRequest = (state: UserSetupQueryState): SearchRequest => {
  const quickStatus = state.quickFilterIds[0];
  const filters = _map(state.filters, (item) => ({
    field: item.field,
    operator: item.operator,
    values: [...item.values],
  }));
  if (quickStatus && quickStatus !== 'ALL') {
    filters.push({ field: UserField.STATUS, operator: 'equals', values: [quickStatus] });
  }
  return {
    searchText: state.searchText.trim(),
    filters,
    sort: _map(state.sorting.slice(0, 1), (rule) => ({ field: rule.field, direction: rule.direction })),
    pageInfo: { rows: state.pageSize, cursor: state.cursors[state.cursors.length - 1] },
  };
};
