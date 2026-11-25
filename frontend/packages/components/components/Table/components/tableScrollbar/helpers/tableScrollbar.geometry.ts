export const getScrollbarGeometry = (visible: number, content: number, offset: number) => {
  const maximum = Math.max(0, content - visible);
  if (!maximum || visible <= 0) return { maximum: 0, size: 100, position: 0 };
  const size = Math.min(100, Math.max(visible / content * 100, 24 / visible * 100));
  const position = Math.min(1, Math.max(0, offset / maximum)) * (100 - size);
  return { maximum, size, position };
};

export const getScrollbarOffset = (pointer: number, trackLength: number, thumbSize: number, grab: number, maximum: number): number => {
  const travel = trackLength * (1 - thumbSize / 100);
  if (travel <= 0) return 0;
  return Math.min(maximum, Math.max(0, (pointer - grab) / travel * maximum));
};
