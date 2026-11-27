import type { FiltersProps } from '../../../types/filters.props';
import type { useFilters } from '../../../hooks/useFilters';

export type MobileFilterListProps = Pick<FiltersProps, 'fields'> & {
  readonly controls: ReturnType<typeof useFilters>;
};
