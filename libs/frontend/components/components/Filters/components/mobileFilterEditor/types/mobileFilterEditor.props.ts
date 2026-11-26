import type { FilterCondition } from '../../../types/filters.props';
import type { useFilters } from '../../../hooks/useFilters';

export interface MobileFilterEditorProps {
  readonly editing: FilterCondition;
  readonly controls: ReturnType<typeof useFilters>;
}
