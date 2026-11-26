export const reorderDragList = <T>(
  items: readonly T[],
  from: number,
  to: number,
  canMove?: (item: T, target: T) => boolean,
): readonly T[] => {
  const item = items[from];
  if (item === undefined || from === to || to < 0 || to >= items.length) return items;
  const start = Math.min(from, to);
  const end = Math.max(from, to);
  if (canMove) {
    for (const target of items.slice(start, end + 1)) {
      if (target !== item && !canMove(item, target)) return items;
    }
  }
  const result = [...items];
  result.splice(from, 1);
  result.splice(to, 0, item);
  return result;
};
