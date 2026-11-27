/** Single source of query defaults for the SPA. */
export const SHELL_QUERY_OPTIONS = {
  defaultOptions: { queries: { staleTime: 30_000, retry: 1, refetchOnWindowFocus: false } },
} as const;
