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
  setSortable(sortable: boolean): Column<T> {
    return new Column({ ...this.definition, sortable });
  }
  setSortField(sortField: string): Column<T> {
    return new Column({ ...this.definition, sortField });
  }
  setWidth(width: number): Column<T> {
    return new Column({ ...this.definition, width });
  }
  setHideable(hideable: boolean): Column<T> {
    return new Column({ ...this.definition, hideable });
  }
  setVisible(visible: boolean): Column<T> {
    return new Column({ ...this.definition, visible });
  }
  build(): TableColumn<T> {
    const { id, header, accessor, width } = this.definition;
    if (!id || !header || !accessor) throw new Error('Column requires id, header and accessor');
    if (width !== undefined && (!Number.isFinite(width) || width <= 0)) throw new Error('Invalid column width');
    return Object.freeze({ ...this.definition, id, header, accessor });
  }
}
