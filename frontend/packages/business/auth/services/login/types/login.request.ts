/** Exact request fields accepted by the identity login endpoint. */
export interface LoginRequest {
  identifier: string;
  password: string;
}
