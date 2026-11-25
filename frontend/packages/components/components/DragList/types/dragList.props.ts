import type { ListProps } from '../../List';

export interface DragListProps<T> extends ListProps<T> {
  readonly onReorder: (items: readonly T[]) => void;
  readonly canMove?: (item: T, target: T) => boolean;
  readonly disabled?: boolean;
  readonly dataTestId?: string;
}
