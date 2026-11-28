import type { HttpMethod } from '../constants/types/api.method';

export interface ApiOptions {
  readonly method?: HttpMethod;
  readonly body?: unknown;
  readonly signal?: AbortSignal;
  readonly headers?: Record<string, string>;
}
