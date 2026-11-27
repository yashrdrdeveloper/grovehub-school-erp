import type { ReactElement } from 'react';
import type { ForgotPasswordViewModel } from '../../../../types/forgotPassword.viewModel';

export type RenderForgotPasswordNotice = (view: ForgotPasswordViewModel) => ReactElement;
