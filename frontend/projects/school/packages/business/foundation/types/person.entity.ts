export interface Person {
  readonly gid: string;
  readonly firstName: string;
  readonly middleName: string | null;
  readonly lastName: string;
  readonly status: string;
  readonly version: number;
}
