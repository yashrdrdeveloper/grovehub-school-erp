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
