import { _map } from '../map';

describe('_map', () => {
  it('transforms values in order without mutating input', () => {
    const entries = Object.freeze([1, 2, 3]);
    expect(_map(entries, value => value * 2)).toEqual([2, 4, 6]);
    expect(entries).toEqual([1, 2, 3]);
  });
});
