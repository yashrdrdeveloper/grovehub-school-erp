import { _find } from '../find';

describe('_find', () => {
  it('returns the matching entry without modifying the collection', () => {
    const values = Object.freeze([{ id: 1 }, { id: 2 }]);
    expect(_find(values, item => item.id === 2)).toBe(values[1]);
    expect(values).toHaveLength(2);
  });
  it('returns undefined when no entry matches', () => {
    expect(_find([1, 2], item => item === 4)).toBeUndefined();
  });
});
