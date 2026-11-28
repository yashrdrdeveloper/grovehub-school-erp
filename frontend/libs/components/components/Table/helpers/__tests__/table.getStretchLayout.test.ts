import { TablePin } from '../../constants/table.type';
import type { TableColumnLayout } from '../../types/table.row';
import { getTableStretchLayout } from '../table.getStretchLayout';

const layout: readonly TableColumnLayout<string>[] = [
  {
    column: { id: 'user', header: 'User', accessor: (row) => row, fillRemainingWidth: true },
    width: 300,
    pin: TablePin.LEFT,
    offset: 0,
  },
  {
    column: { id: 'status', header: 'Status', accessor: (row) => row },
    width: 160,
    pin: TablePin.LEFT,
    offset: 300,
  },
  {
    column: { id: 'date', header: 'Date', accessor: (row) => row },
    width: 220,
    pin: TablePin.NONE,
    offset: 0,
  },
];

describe('getTableStretchLayout', () => {
  it('fills spare space and recomputes pinned offsets without mutating columns', () => {
    const result = getTableStretchLayout(layout, 1000, 56, {});
    expect(result.map((entry) => entry.width)).toEqual([564, 160, 220]);
    expect(result[1]?.offset).toBe(564);
    expect(layout[0]?.width).toBe(300);
  });

  it('preserves manual widths and narrow viewport overflow', () => {
    expect(getTableStretchLayout(layout, 1000, 56, { user: 300 })).toBe(layout);
    expect(getTableStretchLayout(layout, 500, 56, {})).toBe(layout);
  });
});
