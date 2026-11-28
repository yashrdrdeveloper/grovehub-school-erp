import type { FormEvent } from 'react';
import type { Operator } from '../../../../../business/api/types';

export interface LoginPageProps {
  readonly onLogin: (operator: Operator) => void;
}
export interface LoginFormProps {
  readonly busy: boolean;
  readonly error: string;
  readonly onSubmit: (event: FormEvent<HTMLFormElement>) => Promise<void>;
}
