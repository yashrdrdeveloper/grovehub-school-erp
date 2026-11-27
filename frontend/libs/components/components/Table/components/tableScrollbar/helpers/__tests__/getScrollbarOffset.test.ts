import { getScrollbarOffset } from '../tableScrollbar.geometry';

describe('getScrollbarOffset', () => {
  it('maps pointer movement back to the native scroll position and clamps bounds', () => {
    expect(getScrollbarOffset(75, 200, 50, 25, 200)).toBe(100);
    expect(getScrollbarOffset(-10, 200, 50, 25, 200)).toBe(0);
    expect(getScrollbarOffset(200, 200, 50, 25, 200)).toBe(200);
  });
});
