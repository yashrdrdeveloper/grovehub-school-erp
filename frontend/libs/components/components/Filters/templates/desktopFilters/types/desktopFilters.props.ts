import type { FiltersProps } from '../../../types/filters.props';
import type { useFilters } from '../../../hooks/useFilters';

export type DesktopFiltersProps = Pick<FiltersProps, 'fields'> & {
  readonly controls: ReturnType<typeof useFilters>;
};
