import { getTableFillWidth } from '../table.getFillWidth';

describe('getTableFillWidth', () => {
  it('fills remaining width without shrinking data columns', () => {
    expect(getTableFillWidth(1000, 700, 56)).toBe(244);
    expect(getTableFillWidth(500, 700, 56)).toBe(0);
  });
});
