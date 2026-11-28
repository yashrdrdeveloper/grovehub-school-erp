import { isEqual } from '../isEqual';

describe('isEqual', () => {
  it('deeply compares equal values through the real pure implementation', () => {
    expect(isEqual({ members: ['one'] }, { members: ['one'] })).toBe(true);
  });
  it('rejects nested differences and mismatched types', () => {
    expect(isEqual({ members: ['one'] }, { members: ['two'] })).toBe(false);
    expect(isEqual([1], { 0: 1 })).toBe(false);
  });
});
