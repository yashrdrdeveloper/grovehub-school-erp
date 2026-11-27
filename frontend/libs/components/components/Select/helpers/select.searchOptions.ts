import { _filter } from '@Base/utilities/filter';

import type { SelectOption, SelectValue } from '../types/select.props';

export const searchSelectOptions = <T extends SelectValue>(
  options: readonly SelectOption<T>[],
  search: string,
): readonly SelectOption<T>[] => {
  const query = search.trim().toLowerCase();
  return _filter(options, (option) => option.label.toLowerCase().includes(query));
};
