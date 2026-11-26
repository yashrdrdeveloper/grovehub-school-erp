import type { MobileTableProps } from '../../../types/mobileTable.props';
import type { TableEntry } from '../../../../../types/table.row';

export interface MobileTableRowProps<T> {
  readonly props: MobileTableProps<T>;
  readonly entry: Extract<TableEntry<T>, { kind: 'ROW' }>;
}
