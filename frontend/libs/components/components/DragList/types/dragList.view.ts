import type { DragListProps } from './dragList.props';
import type { useDragList } from '../hooks/useDragList';

export type DragListViewProps<T> = DragListProps<T> & {
  readonly controls: ReturnType<typeof useDragList<T>>;
};
