import type { TableColumn } from '../../types/table.column';

export class Column<T> {
  constructor(private readonly definition: Readonly<Partial<TableColumn<T>>> = {}) {}

  setId(id: string): Column<T> {
    return new Column({ ...this.definition, id });
  }

  setHeader(header: string): Column<T> {
    return new Column({ ...this.definition, header });
  }

  setAccessor(accessor: keyof T | ((row: T) => unknown)): Column<T> {
    if (typeof accessor === 'function') return new Column({ ...this.definition, accessor });
    return new Column({ ...this.definition, accessor: (row) => row[accessor] });
  }

  setRenderer(renderer: NonNullable<TableColumn<T>['renderer']>): Column<T> {
    return new Column({ ...this.definition, renderer });
  }

  setEditRenderer(editRenderer: NonNullable<TableColumn<T>['editRenderer']>): Column<T> {
    return new Column({ ...this.definition, editRenderer });
  }

  setComparator(compare: NonNullable<TableColumn<T>['compare']>): Column<T> {
    return new Column({ ...this.definition, compare });
  }

  setSortable(sortable: boolean): Column<T> {
    return new Column({ ...this.definition, sortable });
  }

  setSortField(sortField: string): Column<T> {
    return new Column({ ...this.definition, sortField });
  }

  setWidth(width: number): Column<T> {
    return new Column({ ...this.definition, width });
  }

  setMinWidth(minWidth: number): Column<T> {
    return new Column({ ...this.definition, minWidth });
  }

  setMaxWidth(maxWidth: number): Column<T> {
    return new Column({ ...this.definition, maxWidth });
  }

  setPin(pin: TableColumn<T>['pin']): Column<T> {
    return new Column({ ...this.definition, pin });
  }

  setAlignment(alignment: TableColumn<T>['alignment']): Column<T> {
    return new Column({ ...this.definition, alignment });
  }

  setHideable(hideable: boolean): Column<T> {
    return new Column({ ...this.definition, hideable });
  }

  setVisible(visible: boolean): Column<T> {
    return new Column({ ...this.definition, visible });
  }

  setResizable(resizable: boolean): Column<T> {
    return new Column({ ...this.definition, resizable });
  }

  setReorderable(reorderable: boolean): Column<T> {
    return new Column({ ...this.definition, reorderable });
  }

  setMobileVisible(mobileVisible: boolean): Column<T> {
    return new Column({ ...this.definition, mobileVisible });
  }

  setSearchable(searchable: boolean): Column<T> {
    return new Column({ ...this.definition, searchable });
  }

  setExportable(exportable: boolean): Column<T> {
    return new Column({ ...this.definition, exportable });
  }

  build(): TableColumn<T> {
    const { id, header, accessor, width, minWidth, maxWidth } = this.definition;
    if (!id || !header || !accessor) throw new Error('Column requires id, header and accessor');
    for (const size of [width, minWidth, maxWidth]) {
      if (size !== undefined && (!Number.isFinite(size) || size <= 0)) {
        throw new Error('Invalid column width');
      }
    }
    if (minWidth !== undefined && maxWidth !== undefined && minWidth > maxWidth) {
      throw new Error('Minimum column width exceeds maximum');
    }
    return Object.freeze({ ...this.definition, id, header, accessor });
  }
}
