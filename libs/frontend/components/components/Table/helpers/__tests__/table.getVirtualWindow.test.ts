import { getVirtualWindow } from '../table.getVirtualWindow';

describe('getVirtualWindow', () => {
  it('keeps visible rows and overscan without rendering the entire dataset', () => {
    expect(getVirtualWindow(1000, 500, 200, 50, 2)).toEqual({
      start: 8,
      end: 16,
      before: 400,
      after: 49200,
    });
  });

  it('clamps to the dataset bounds', () => {
    expect(getVirtualWindow(2, 500, 200, 50, 2)).toEqual({
      start: 0,
      end: 2,
      before: 0,
      after: 0,
    });
  });
});
