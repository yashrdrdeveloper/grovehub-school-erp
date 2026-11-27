import type { Page } from './types';

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;
  constructor(status: number, code: string) { super(code); this.status = status; this.code = code; }
}

export function csrfFromCookie(cookie: string): string {
  const prefix = '__Host-gh_internal_csrf=';
  const values = cookie.split(';').map((part) => part.trim()).filter((part) => part.startsWith(prefix));
  if (values.length !== 1) return '';
  return values[0]?.slice(prefix.length) ?? '';
}

export async function request<T>(path: string, method = 'GET', body?: unknown, signal?: AbortSignal): Promise<T> {
  if (!/^\/internal\/(public|protected)\/[a-zA-Z0-9_/-]+$/.test(path)) throw new ApiError(400, 'INVALID_PATH');
  const headers: Record<string, string> = { Accept: 'application/json' };
  if (method !== 'GET') {
    const csrf = csrfFromCookie(document.cookie);
    if (csrf) headers['X-CSRF-Token'] = csrf;
    if (body !== undefined) headers['Content-Type'] = 'application/json';
  }
  const options: RequestInit = { method, headers, credentials: 'same-origin', redirect: 'error', cache: 'no-store', signal };
  if (body !== undefined) options.body = JSON.stringify(body);
  const response = await fetch(path, options);
  if (response.status === 401 && path.startsWith('/internal/protected/')) window.dispatchEvent(new Event('internal-session-expired'));
  const result = await response.json().catch(() => null) as { success?: boolean; data?: T; errorCode?: string } | null;
  if (!response.ok || result?.success !== true) throw new ApiError(response.status, result?.errorCode ?? 'UNAVAILABLE');
  return result.data as T;
}

export function search<T>(resource: string, pageNumber: number, rowSize: number, signal?: AbortSignal): Promise<Page<T>> {
  return request<Page<T>>('/internal/protected/' + resource + '/search', 'POST', { pageNumber, rowSize }, signal);
}

export function errorText(error: unknown): string {
  if (error instanceof ApiError) {
    if (error.status === 401) return 'Invalid credentials or expired session.';
    if (error.status === 403) return 'You do not have permission for this action.';
    if (error.status === 409) return 'The record has changed. Refresh and try again.';
    if (error.status === 429) return 'Too many attempts. Try again later.';
    if (error.code === 'EXECUTION_UNAVAILABLE') return 'This template has no configured execution adapter.';
    if (error.status === 400) return 'Check the input and try again.';
    if (error.status === 404) return 'Record not found.';
  }
  return 'Unable to complete the request. Try again.';
}
