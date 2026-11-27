import { getAuditSearchRequest } from '../audit.getSearchRequest';

describe('getAuditSearchRequest', () => {
  it('trims search text and preserves pagination', () => {
    expect(getAuditSearchRequest('  test  ', [], [], 2, 25)).toEqual({
      searchText: 'test',
      filters: [],
      sort: [],
      pageInfo: { pageNumber: 2, rowSize: 25 },
    });
  });
});
