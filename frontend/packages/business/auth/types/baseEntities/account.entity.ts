/** Backend-owned account identity; do not redefine fields in UI pages. */
export interface AccountEntity {
  accountId: string;
  personGid: string;
  username: string | null;
  email: string;
  status: string;
  version: number;
  createdAt: string;
  updatedAt: string;
}
