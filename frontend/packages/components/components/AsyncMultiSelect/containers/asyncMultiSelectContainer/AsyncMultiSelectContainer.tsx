import { useSelection } from '../../../Select/hooks/useSelection';
import { useAsyncOptions } from '../../../AsyncSelect/hooks/useAsyncOptions';

import { AsyncMultiSelectView } from '../../AsyncMultiSelect';
import type { SelectValue } from '../../../Select';
import type { AsyncMultiSelectProps } from '../../types/asyncMultiSelect.props';

export const AsyncMultiSelectContainer = <T extends SelectValue>({
  value,
  onChange,
  cacheKey,
  loadOptions,
  initialOptions,
  resolveOptions,
  pageSize,
  debounceMs,
  maxSelections,
  ...props
}: AsyncMultiSelectProps<T>) => {
  const controls = useSelection({ value, onChange, multiple: true, maxSelections });
  const source = useAsyncOptions({ cacheKey, loadOptions, initialOptions, resolveOptions, pageSize, debounceMs }, controls.open, value);
  const getSelected = () => {
    if (controls.open) return controls.selected;
    return value;
  };

  return (
    <AsyncMultiSelectView
      {...props}
      {...controls}
      selected={getSelected()}
      multiple
      options={source.options}
      displayOptions={source.displayOptions}
      loading={source.loading}
      error={source.error || props.error}
      hasNext={source.hasNext}
      onLoadMore={source.onLoadMore}
      onRetry={source.onRetry}
      onOpenChange={(open) => {
        if (open) source.search('');
        controls.onOpenChange(open);
      }}
      onSearch={(search) => {
        controls.onSearch(search);
        source.search(search);
      }}
    />
  );
};
