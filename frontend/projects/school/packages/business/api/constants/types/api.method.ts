import type { HTTP_METHOD } from '../api.method';

export type HttpMethod = (typeof HTTP_METHOD)[keyof typeof HTTP_METHOD];
