import { TablePin } from '../../constants/table.type';
import type { TableColumnLayout } from '../../types/table.row';
import { getTablePinnedBoundary } from '../table.getPinnedBoundary';

describe('getTablePinnedBoundary', () => {
  it('marks only the edge of each pinned group facing scrolling columns', () => {
    const layout: TableColumnLayout<string>[] = [TablePin.LEFT, TablePin.LEFT, TablePin.NONE, TablePin.RIGHT, TablePin.RIGHT].map((pin, index) => ({
      column: { id: String(index), header: String(index), accessor: (row) => row },
      width: 100,
      offset: 0,
      pin,
    }));
    expect(layout.map((_, index) => getTablePinnedBoundary(layout, index))).toEqual([TablePin.NONE, TablePin.LEFT, TablePin.NONE, TablePin.RIGHT, TablePin.NONE]);
  });
});
