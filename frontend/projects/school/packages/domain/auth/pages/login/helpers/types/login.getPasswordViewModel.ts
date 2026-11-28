import type { LucideIcon } from '@Components/icons';
export interface LoginPasswordViewModel {
  readonly Icon: LucideIcon;
  readonly label: string;
  readonly type: 'text' | 'password';
}
export type GetLoginPasswordViewModel = (visible: boolean) => LoginPasswordViewModel;
