export interface ApiSuccess<T> {
  success: true;
  data: T;
}

export interface ApiErrorBody {
  success: false;
  errorCode: string;
  errorMessage: string;
  httpErrorCode: number;
  errorMetaData?: Record<string, unknown>;
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
