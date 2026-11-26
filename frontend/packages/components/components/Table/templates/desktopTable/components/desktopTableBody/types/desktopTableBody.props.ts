import type { DesktopTableProps } from '../../../types/desktopTable.props';
import type { TableColumnLayout } from '../../../../../types/table.row';

import type { CSSProperties } from 'react';
import type { TableEntry } from '../../../../../types/table.row';
import type { getVirtualWindow } from '../../../../../helpers/table.getVirtualWindow';

export interface DesktopTableBodyProps<T> {
  readonly props: DesktopTableProps<T>;
  readonly layout: readonly TableColumnLayout<T>[];
  readonly selection: boolean;
  readonly rowActions: boolean;
  readonly fillBefore: (index: number) => boolean;
  readonly fillAtEnd: boolean;
  readonly entries: readonly TableEntry<T>[];
  readonly span: number;
  readonly virtual: boolean;
  readonly range: ReturnType<typeof getVirtualWindow>;
  readonly rowStyle: () => CSSProperties | undefined;
}
