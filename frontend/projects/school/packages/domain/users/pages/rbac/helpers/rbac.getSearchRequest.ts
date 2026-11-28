import { _map } from '@Base/utilities/map';
import { _trim } from '@Base/utilities/trim';
import type { FilterCondition } from '@Components/components/Filters';
import type { SortCondition } from '@Components/components/Sort';
import type { SearchRequest } from '@Business/accounts/types';

export const getRbacSearchRequest = (
  searchText: string,
  filters: readonly FilterCondition[],
  sorting: readonly SortCondition[],
  pageNumber: number,
  rowSize: number,
): SearchRequest => ({
  searchText: _trim(searchText),
  filters: _map(filters, (filter) => ({
    field: filter.field,
    operator: 'EQUALS',
    values: [...filter.values],
  })),
  sort: _map(sorting, (sort) => ({ field: sort.field, direction: sort.direction })),
  pageInfo: { pageNumber, rowSize },
});
