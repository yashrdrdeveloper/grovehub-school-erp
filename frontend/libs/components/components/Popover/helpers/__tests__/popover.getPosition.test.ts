import { getPopoverPosition } from '../popover.getPosition';

describe('getPopoverPosition', () => {
  it('places a footer dropdown immediately above its anchor using its measured height', () => {
    const result = getPopoverPosition({
      anchor: { left: 80, top: 700, bottom: 732, height: 32 },
      panelHeight: 128,
      viewportWidth: 1000,
      viewportHeight: 760,
      width: 88,
    });
    expect(result.top).toBe(568);
    expect(result.left).toBe(80);
  });
  it('places a dropdown below its anchor when space permits', () => {
    const result = getPopoverPosition({
      anchor: { left: 80, top: 100, bottom: 132, height: 32 },
      panelHeight: 128,
      viewportWidth: 1000,
      viewportHeight: 760,
      width: 180,
    });
    expect(result.top).toBe(136);
  });
});
