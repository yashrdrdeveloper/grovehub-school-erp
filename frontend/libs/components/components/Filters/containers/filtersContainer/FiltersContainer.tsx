import { FiltersView } from '../../Filters';
import { useFilters } from '../../hooks/useFilters';
import type { FiltersProps } from '../../types/filters.props';

export const FiltersContainer = (props: FiltersProps) => {
  const controls = useFilters(props);
  return <FiltersView {...props} controls={controls} />;
};
