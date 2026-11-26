import type { UserFormProps } from '../../../types/userForm.props';

export type UserAccountActionsProps = Pick<UserFormProps, 'row' | 'submitting' | 'onAccountAction'>;
