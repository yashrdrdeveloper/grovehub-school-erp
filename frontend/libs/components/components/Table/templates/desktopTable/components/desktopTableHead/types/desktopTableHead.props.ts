import type { DesktopTableProps } from '../../../types/desktopTable.props';
import type { TableColumnLayout } from '../../../../../types/table.row';

export interface DesktopTableHeadProps<T> {
  readonly props: DesktopTableProps<T>;
  readonly layout: readonly TableColumnLayout<T>[];
  readonly selection: boolean;
  readonly rowActions: boolean;
  readonly fillBefore: (index: number) => boolean;
  readonly fillAtEnd: boolean;
  readonly fill: number;
}
