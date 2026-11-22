import type { LoginProps } from '../../../types';
import type { LoginFormModel } from '../../useLoginForm';

export type LoginActions = Pick<LoginProps,
  'onEmailChange' | 'onPasswordChange' | 'onTogglePassword'>;

export type UseLoginActions = (model: LoginFormModel) => LoginActions;
