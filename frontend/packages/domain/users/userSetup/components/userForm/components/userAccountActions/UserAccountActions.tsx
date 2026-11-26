import { UserStatus } from '@Business/domain/users';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';

import { UserAction } from '../../../../constants/userSetup.actions';
import type { UserAccountActionsProps } from './types/userAccountActions.props';

export const UserAccountActions = ({ row, submitting, onAccountAction }: UserAccountActionsProps) => (
  <>
    {row?.status === UserStatus.SUSPENDED && (
      <Button
        variant={BUTTON_VARIANT.OUTLINE}
        disabled={submitting}
        onClick={() => onAccountAction?.(UserAction.ACTIVATE, row)}
      >
        {__('Activate account')}
      </Button>
    )}
    {row?.status === UserStatus.PENDING_SETUP && (
      <Button
        variant={BUTTON_VARIANT.OUTLINE}
        disabled={submitting}
        onClick={() => onAccountAction?.(UserAction.RESEND, row)}
      >
        {__('Resend invitation')}
      </Button>
    )}
  </>
);
