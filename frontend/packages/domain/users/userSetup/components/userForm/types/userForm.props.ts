import type { Account, CreateAccount } from '@Business/domain/users';

export interface UserFormProps {
  readonly row?: Account;
  readonly onClose: () => void;
  readonly onCreate: (input: CreateAccount, key: string) => void;
  readonly onEdit: (row: Account, username: string | null) => void;
  readonly submitting: boolean;
  readonly error: boolean;
}
