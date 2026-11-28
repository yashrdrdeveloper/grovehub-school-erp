import { Column } from '@Components/components/Table';
import type { RecordColumn } from '../types/pagedRecords.props';

export const getRecordColumns = <T>(columns: readonly RecordColumn<T>[]) =>
  columns.map((column) =>
    new Column<T>()
      .setId(column.id)
      .setHeader(column.title)
      .setAccessor(column.cell)
      .setRenderer(column.cell)
      .setMinWidth(110)
      .setWidth(190)
      .setSortable(false)
      .build(),
  );
