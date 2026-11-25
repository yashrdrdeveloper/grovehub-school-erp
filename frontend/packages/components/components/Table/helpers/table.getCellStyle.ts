import type { CSSProperties } from 'react';
import { TABLE_DIMENSIONS } from '../constants/table.dimensions';
import { TableAlignment, TablePin } from '../constants/table.type';
import type { TableColumnLayout } from '../types/table.row';

export const getTableCellStyle = <T,>(entry: TableColumnLayout<T>, selection: boolean, actions: boolean): CSSProperties => {
  const style: CSSProperties = { width: entry.width, minWidth: entry.width, maxWidth: entry.width };
  if (entry.column.alignment === TableAlignment.CENTER) style.textAlign = 'center';
  if (entry.column.alignment === TableAlignment.END) style.textAlign = 'end';
  if (entry.pin === TablePin.LEFT) {
    style.position = 'sticky';
    style.left = entry.offset + Number(selection) * TABLE_DIMENSIONS.selection;
    style.zIndex = 1;
  }
  if (entry.pin === TablePin.RIGHT) {
    style.position = 'sticky';
    style.right = entry.offset + Number(actions) * TABLE_DIMENSIONS.actions;
    style.zIndex = 1;
  }
  return style;
};
