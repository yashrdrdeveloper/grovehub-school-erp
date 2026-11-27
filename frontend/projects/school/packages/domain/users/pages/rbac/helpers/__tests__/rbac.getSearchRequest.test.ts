import { getRbacSearchRequest } from '../rbac.getSearchRequest';

describe('getRbacSearchRequest', () => {
  it('trims search text and preserves pagination', () => {
    expect(getRbacSearchRequest('  test  ', [], [], 2, 25)).toEqual({ searchText: 'test', filters: [], sort: [], pageInfo: { pageNumber: 2, rowSize: 25 } });
  });
});
