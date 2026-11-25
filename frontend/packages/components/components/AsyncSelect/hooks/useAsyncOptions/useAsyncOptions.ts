import { useState } from 'react';
import { useInfiniteQuery } from '@tanstack/react-query';
import { _map } from '@Base/utilities/map';
import { useDebounce } from '@Components/hooks/useDebounce';

import type { SelectOption, SelectValue } from '../../../Select';
import type { AsyncSelectSource } from '../../types/asyncSelect.props';

export const useAsyncOptions = <T extends SelectValue>(
  source: AsyncSelectSource<T>,
  enabled: boolean,
) => {
  const [query, setQuery] = useState('');
  const { execute, cancel } = useDebounce(setQuery, source.debounceMs ?? 350);
  const result = useInfiniteQuery({
    queryKey: ['shared-async-select', ...source.cacheKey, query, source.pageSize],
    enabled,
    initialPageParam: undefined as string | undefined,
    queryFn: ({ pageParam, signal }) => source.loadOptions({
      searchText: query,
      cursor: pageParam,
      limit: source.pageSize ?? 25,
      signal,
    }),
    getNextPageParam: (page) => {
      if (page.hasNext) return page.nextCursor;
      return undefined;
    },
  });
  const knownOptions = new Map<T, SelectOption<T>>();
  for (const option of source.initialOptions ?? []) knownOptions.set(option.value, option);
  for (const page of result.data?.pages ?? []) {
    for (const option of page.options) knownOptions.set(option.value, option);
  }
  const search = (next: string) => {
    if (!next) {
      cancel();
      setQuery('');
      return;
    }
    execute(next);
  };
  return {
    options: _map(Array.from(knownOptions.values()), (option) => option),
    search,
    loading: result.isFetching,
    error: result.isError && __('Unable to load options'),
    hasNext: result.hasNextPage,
    onLoadMore: () => {
      if (result.hasNextPage && !result.isFetching) void result.fetchNextPage();
    },
    onRetry: () => { void result.refetch(); },
  };
};
