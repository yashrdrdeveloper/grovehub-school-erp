import { FilterOperator } from '@Components/components/Filters';
import { SortDirection } from '@Components/components/Sort';

import { getUserSetupSearchRequest } from '../userSetup.getSearchRequest';

describe('getUserSetupSearchRequest', () => {
  it('serializes inclusive calendar ranges and preserves ordered sort rules', () => {
    const request = getUserSetupSearchRequest({
      searchText: ' test ',
      filters: [
        { id: '1', field: 'status', operator: FilterOperator.EQUALS, values: ['ACTIVE'] },
        { id: '2', field: 'createdAt', operator: FilterOperator.BETWEEN, values: ['2026-11-01', '2026-11-03'] },
        { id: '3', field: 'updatedAt', operator: FilterOperator.BETWEEN, values: ['2026-11-05', '2026-11-05'] },
      ],
      sorting: [
        { field: 'status', direction: SortDirection.ASC },
        { field: 'createdAt', direction: SortDirection.DESC },
      ],
      quickFilterIds: ['ALL'],
      pageSize: 25,
      cursors: [undefined],
    });
    expect(request.filters).toEqual([
      { field: 'status', operator: 'EQUALS', values: ['ACTIVE'] },
      { field: 'createdAt', operator: 'GTE', values: [new Date(2026, 10, 1).toISOString()] },
      { field: 'createdAt', operator: 'LT', values: [new Date(2026, 10, 4).toISOString()] },
      { field: 'updatedAt', operator: 'GTE', values: [new Date(2026, 10, 5).toISOString()] },
      { field: 'updatedAt', operator: 'LT', values: [new Date(2026, 10, 6).toISOString()] },
    ]);
    expect(request.sort).toEqual([
      { field: 'status', direction: 'ASC' },
      { field: 'createdAt', direction: 'DESC' },
    ]);
  });
});
