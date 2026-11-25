import { SelectView } from '../../Select';
import { useSelection } from '../../hooks/useSelection';
import { searchSelectOptions } from '../../helpers/select.searchOptions';

import type { SelectProps, SelectValue } from '../../types/select.props';

export const SelectContainer = <T extends SelectValue>({
  value,
  onChange,
  options,
  ...props
}: SelectProps<T>) => {
  const selected: T[] = [];
  if (value !== null) selected.push(value);
  const controls = useSelection<T>({
    value: selected,
    multiple: false,
    onChange: (values) => onChange(values[0] ?? null),
  });
  return (
    <SelectView
      {...props}
      {...controls}
      selected={selected}
      multiple={false}
      displayOptions={options}
      options={searchSelectOptions(options, controls.search)}
    />
  );
};
