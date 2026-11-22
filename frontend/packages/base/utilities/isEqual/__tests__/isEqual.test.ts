import { isEqual } from '../isEqual';

describe('isEqual', () => {
  it('deeply compares value structure', () => {
    expect(isEqual({ members: ['one', 'two'] }, { members: ['one', 'two'] })).toBe(true);
    expect(isEqual({ members: ['one'] }, { members: ['two'] })).toBe(false);
  });
});
