import { useSelection } from '../../../Select/hooks/useSelection';

import { AsyncSelectView } from '../../AsyncSelect';
import { useAsyncOptions } from '../../hooks/useAsyncOptions';

import type { SelectValue } from '../../../Select';
import type { AsyncSelectProps } from '../../types/asyncSelect.props';

export const AsyncSelectContainer = <T extends SelectValue>({
  value,
  onChange,
  cacheKey,
  loadOptions,
  initialOptions,
  resolveOptions,
  pageSize,
  debounceMs,
  ...props
}: AsyncSelectProps<T>) => {
  const selected: T[] = [];
  if (value !== null) selected.push(value);
  const controls = useSelection<T>({
    value: selected,
    multiple: false,
    onChange: (values) => onChange(values[0] ?? null),
  });
  const source = useAsyncOptions(
    { cacheKey, loadOptions, initialOptions, resolveOptions, pageSize, debounceMs },
    controls.open,
    selected,
  );

  return (
    <AsyncSelectView
      {...props}
      {...controls}
      selected={selected}
      multiple={false}
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
