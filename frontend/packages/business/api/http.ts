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
  const headers: Record<string, string> = {};
  if (options.body !== undefined) headers['Content-Type'] = 'application/json';
  if (method !== 'GET' && csrf) headers['X-CSRF-Token'] = csrf;
  if (options.headers) Object.assign(headers, options.headers);
  const request: RequestInit = {
    method,
    credentials: 'include',
    cache: 'no-store',
    signal: options.signal,
    headers,
  };
  if (options.body !== undefined) request.body = JSON.stringify(options.body);
  const response = await fetch(path, request);

  if (response.status === 204) {
    return undefined as T;
  }

  const payload: unknown = await response.json().catch(() => undefined);
  if (!response.ok) {
    const failure = payload as ApiErrorBody | undefined;
    throw new ApiError(
      response.status,
      failure?.errorCode ?? 'HTTP_ERROR',
      failure?.errorMessage ?? `Request failed (${response.status})`,
      failure?.errorMetaData,
    );
  }

  const envelope = payload as ApiSuccess<T> | undefined;
  if (!envelope || envelope.success !== true) {
    throw new ApiError(response.status, 'INVALID_RESPONSE', 'Unexpected API response envelope');
  }
  return envelope.data;
};
