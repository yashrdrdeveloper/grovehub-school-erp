import { reorderDragList } from '../dragList.reorder';

describe('reorderDragList', () => {
  it('moves items without mutating the source', () => {
    const items = ['a', 'b', 'c'];
    expect(reorderDragList(items, 0, 2)).toEqual(['b', 'c', 'a']);
    expect(items).toEqual(['a', 'b', 'c']);
  });

  it('does not cross a non-reorderable item', () => {
    const items = ['a', 'b', 'c'];
    expect(reorderDragList(items, 0, 2, (_, target) => target !== 'b')).toBe(items);
    expect(reorderDragList(items, 0, -1)).toBe(items);
  });
});
