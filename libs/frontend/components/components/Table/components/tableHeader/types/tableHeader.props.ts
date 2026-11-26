import type { TableHeaderProps } from '../../../types/table.props';

export interface TableHeaderViewProps extends TableHeaderProps {
  readonly onAction?: (id: string) => void;
}
