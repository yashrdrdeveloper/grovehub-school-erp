import type { ReactElement } from 'react';
import type { ForgotPasswordProps } from '../../../../types/forgotPassword.props';

export type RenderForgotPasswordSuccess = (
  props: Pick<ForgotPasswordProps, 'formId' | 'state' | 'view' | 'onUseAnotherEmail'>,
) => ReactElement;
