import { useSelection } from '../../../Select/hooks/useSelection';
import { searchSelectOptions } from '../../../Select/helpers/select.searchOptions';

import { MultiSelectView } from '../../MultiSelect';
import type { MultiSelectProps } from '../../types/multiSelect.props';
import type { SelectValue } from '../../../Select';

export const MultiSelectContainer = <T extends SelectValue>({
  value,
  onChange,
  options,
  maxSelections,
  ...props
}: MultiSelectProps<T>) => {
  const controls = useSelection({ value, onChange, multiple: true, maxSelections });
  const getSelected = () => {
    if (controls.open) return controls.selected;
    return value;
  };

  return (
    <MultiSelectView
      {...props}
      {...controls}
      selected={getSelected()}
      multiple
      options={searchSelectOptions(options, controls.search)}
    />
  );
};
