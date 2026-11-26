import { SortView } from '../../Sort';
import { useSort } from '../../hooks/useSort';
import type { SortProps } from '../../types/sort.props';

export const SortContainer = (props: SortProps) => {
  const controls = useSort(props);
  return <SortView {...props} controls={controls} />;
};
