/** Structured HTTP error originating from the shared transport boundary. */
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
