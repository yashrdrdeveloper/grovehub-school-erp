import type { AsyncSelectSource } from '../../AsyncSelect';
import type { MultiSelectProps } from '../../MultiSelect';
import type { SelectValue } from '../../Select';

export interface AsyncMultiSelectProps<T extends SelectValue = string>
  extends Omit<MultiSelectProps<T>, 'options'>, AsyncSelectSource<T> {
  readonly dataTestId?: string;
}
