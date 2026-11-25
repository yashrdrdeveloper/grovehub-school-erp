import { DragListView } from '../../DragList';
import { useDragList } from '../../hooks/useDragList';
import type { DragListProps } from '../../types/dragList.props';

export const DragListContainer = <T,>(props: DragListProps<T>) => {
  const controls = useDragList(props);
  return <DragListView {...props} controls={controls} />;
};
