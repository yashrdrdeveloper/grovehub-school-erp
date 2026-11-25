import type { FiltersProps } from '../../../types/filters.props';
import type { useFilters } from '../../../hooks/useFilters';

export type DesktopFiltersProps = Pick<FiltersProps, 'fields' | 'maxFilters'> & { readonly controls: ReturnType<typeof useFilters> };
