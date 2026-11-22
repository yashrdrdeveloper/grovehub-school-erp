import { _filter } from '../filter';

describe('_filter', () => {
  it('returns matching values without mutating the input', () => {
    const original = Object.freeze([1, 2, 3]);
    expect(_filter(original, x => x > 1)).toEqual([2, 3]);
    expect(original).toEqual([1, 2, 3]);
  });
});
