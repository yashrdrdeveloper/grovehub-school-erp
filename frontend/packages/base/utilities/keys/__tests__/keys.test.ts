import { _keys } from '../keys';

describe('_keys', () => {
  it('returns enumerable own keys for a structured object', () => {
    expect(_keys({ email: 'a', password: 'b' })).toEqual(['email', 'password']);
  });
  it('returns no keys for an empty object', () => {
    expect(_keys({})).toEqual([]);
  });
});
