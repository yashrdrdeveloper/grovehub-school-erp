export const getVirtualWindow = (count: number, scrollTop: number, height: number, rowHeight: number, overscan: number) => {
  const size = Math.max(1, rowHeight);
  const total = Math.max(0, count);
  const extra = Math.max(0, overscan);
  const visibleStart = Math.min(total, Math.floor(Math.max(0, scrollTop) / size));
  const start = Math.max(0, visibleStart - extra);
  const end = Math.min(total, Math.max(start, Math.ceil((Math.max(0, scrollTop) + height) / size) + extra));
  return { start, end, before: start * size, after: (total - end) * size };
};
