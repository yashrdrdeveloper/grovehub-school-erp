import { readCSRF } from './cookies';
import { ApiError, type ApiErrorBody, type ApiSuccess } from './types';

type HttpMethod = 'GET' | 'POST' | 'PATCH' | 'DELETE';

export interface ApiOptions {
  method?: HttpMethod;
  body?: unknown;
  signal?: AbortSignal;
  headers?: Record<string, string>;
}

export const apiRequest = async <T>(
  path: string,
  options: ApiOptions = {},
): Promise<T> => {
  const method = options.method ?? 'GET';
  const csrf = readCSRF();
  const response = await fetch(path, {
    method,
    credentials: 'include',
    cache: 'no-store',
    signal: options.signal,
    headers: {
      ...(options.body === undefined ? {} : { 'Content-Type': 'application/json' }),
      ...(method !== 'GET' && csrf ? { 'X-CSRF-Token': csrf } : {}),
      ...options.headers,
    },
    ...(options.body === undefined ? {} : { body: JSON.stringify(options.body) }),
  });

  if (response.status === 204) {
    return undefined as T;
  }

  const payload: unknown = await response.json().catch(() => undefined);
  if (!response.ok) {
    const failure = payload as ApiErrorBody | undefined;
    throw new ApiError(
      response.status,
      failure?.error?.code ?? 'HTTP_ERROR',
      failure?.error?.message ?? `Request failed (${response.status})`,
      failure?.error?.metadata,
    );
  }

  const envelope = payload as ApiSuccess<T> | undefined;
  if (!envelope || envelope.success !== true) {
    throw new ApiError(response.status, 'INVALID_RESPONSE', 'Unexpected API response envelope');
  }
  return envelope.data;
};
