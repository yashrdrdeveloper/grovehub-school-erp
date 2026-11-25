import { getScrollbarGeometry, getScrollbarOffset } from '../tableScrollbar.geometry';

describe('getScrollbarGeometry', () => {
  it('has no thumb when all content fits', () => {
    expect(getScrollbarGeometry(400, 400, 0).maximum).toBe(0);
  });

  it('maps scroll position to proportional thumb travel', () => {
    expect(getScrollbarGeometry(200, 400, 100)).toEqual({ maximum: 200, size: 50, position: 25 });
    expect(getScrollbarGeometry(200, 400, 200).position).toBe(50);
  });

  it('keeps very long content draggable with a minimum thumb', () => {
    expect(getScrollbarGeometry(200, 10000, 0).size).toBe(12);
  });
});

describe('getScrollbarOffset', () => {
  it('maps pointer movement back to the native scroll position and clamps bounds', () => {
    expect(getScrollbarOffset(75, 200, 50, 25, 200)).toBe(100);
    expect(getScrollbarOffset(-10, 200, 50, 25, 200)).toBe(0);
    expect(getScrollbarOffset(200, 200, 50, 25, 200)).toBe(200);
  });
});
