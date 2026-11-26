import { useState } from 'react';
import { useInfiniteQuery, useQuery, useQueryClient } from '@tanstack/react-query';
import type { InfiniteData } from '@tanstack/react-query';
import { useDebounce } from '@Components/hooks/useDebounce';

import type { SelectOption, SelectValue } from '../../../Select';
import type { AsyncOptionsResponse, AsyncSelectSource } from '../../types/asyncSelect.props';

export const useAsyncOptions = <T extends SelectValue>(
  source: AsyncSelectSource<T>,
  enabled: boolean,
  selected: readonly T[],
) => {
  const [query, setQuery] = useState('');
  const client = useQueryClient();
  const { execute, cancel } = useDebounce(setQuery, source.debounceMs ?? 350);
  const scope = ['shared-async-select', ...source.cacheKey];
  const result = useInfiniteQuery({
    queryKey: [...scope, 'options', query, source.pageSize],
    enabled,
    initialPageParam: undefined as string | undefined,
    queryFn: ({ pageParam, signal }) =>
      source.loadOptions({
        searchText: query,
        cursor: pageParam,
        limit: source.pageSize ?? 25,
        signal,
      }),
    getNextPageParam: (page, _pages, _previous, allParams) => {
      if (page.hasNext && page.nextCursor && !allParams.includes(page.nextCursor))
        return page.nextCursor;
      return undefined;
    },
  });
  const known = new Map<T, SelectOption<T>>();
  for (const option of source.initialOptions ?? []) known.set(option.value, option);
  for (const [, cached] of client.getQueriesData<InfiniteData<AsyncOptionsResponse<T>>>({
    queryKey: [...scope, 'options'],
  })) {
    for (const page of cached?.pages ?? []) {
      for (const option of page.options) known.set(option.value, option);
    }
  }
  const missing = selected.filter((value) => !known.has(value));
  const resolved = useQuery({
    queryKey: [...scope, 'selected', missing],
    enabled: Boolean(source.resolveOptions && missing.length),
    queryFn: ({ signal }) => source.resolveOptions?.(missing, signal) ?? Promise.resolve([]),
    staleTime: 300000,
  });
  for (const option of resolved.data ?? []) known.set(option.value, option);
  const options = new Map<T, SelectOption<T>>();
  for (const page of result.data?.pages ?? []) {
    for (const option of page.options) options.set(option.value, option);
  }
  if (!query) {
    for (const option of source.initialOptions ?? []) {
      if (!options.has(option.value)) options.set(option.value, option);
    }
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
    options: Array.from(options.values()),
    displayOptions: Array.from(known.values()),
    search,
    loading: result.isFetching,
    error: result.isError && __('Unable to load options'),
    hasNext: result.hasNextPage,
    onLoadMore: () => {
      if (result.hasNextPage && !result.isFetching) void result.fetchNextPage();
    },
    onRetry: () => {
      void result.refetch();
    },
  };
};
