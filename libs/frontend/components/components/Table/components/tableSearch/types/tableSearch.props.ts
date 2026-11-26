import type { TableSearchProps } from '../../../types/table.props';

export interface TableSearchViewProps extends TableSearchProps {
  readonly onChange: (value: string) => void;
}
