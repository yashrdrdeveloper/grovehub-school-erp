export interface ApiMeta {
  requestId?: string;
  [key: string]: unknown;
}

export interface ApiSuccess<T> {
  success: true;
  data: T;
  meta?: ApiMeta;
}

export interface ApiErrorBody {
  success: false;
  error: {
    code: string;
    message: string;
    metadata?: Record<string, unknown>;
  };
  meta?: ApiMeta;
}

export class ApiError extends Error {
  constructor(
    readonly status: number,
    readonly code: string,
    message: string,
    readonly metadata?: Record<string, unknown>,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}
