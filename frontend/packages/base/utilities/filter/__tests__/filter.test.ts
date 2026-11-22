import { _filter } from '../filter';

describe('_filter', () => {
  it('preserves type-guard narrowing for unknown array values', () => {
    const values: unknown[] = ['valid', 42, 'allowed'];
    const filtered: string[] = _filter(values, (value): value is string => typeof value === 'string');
    expect(filtered).toEqual(['valid', 'allowed']);
  });

  it('returns matching values without mutating the input', () => {
    const original = Object.freeze([1, 2, 3]);
    expect(_filter(original, x => x > 1)).toEqual([2, 3]);
    expect(original).toEqual([1, 2, 3]);
  });
});
