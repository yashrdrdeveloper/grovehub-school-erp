import { getAccountResourcePath } from '../accounts.getResourcePath';

describe('getAccountResourcePath', () => {
  it('encodes an account ID before inserting it into a protected URL', () => {
    expect(getAccountResourcePath('person/one')).toBe('/identity/protected/accounts/person%2Fone');
  });
});
