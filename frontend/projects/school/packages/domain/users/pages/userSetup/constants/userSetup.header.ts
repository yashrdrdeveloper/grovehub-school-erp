import { Plus } from '@Components/icons';
import { UserAction } from './userSetup.actions';

export const USER_SETUP_HEADER = Object.freeze({
  title: __('User Setup'),
  description: __('Manage school accounts and invitations.'),
  rightConfig: [{ id: UserAction.CREATE, label: __('Create user'), icon: Plus }],
});
