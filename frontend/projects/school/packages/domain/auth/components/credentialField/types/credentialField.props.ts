import type { InputHTMLAttributes, ReactElement, ReactNode } from 'react';
import type { LucideIcon } from '@Components/icons';

export interface CredentialFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  readonly id: string;
  readonly label: string;
  readonly Icon: LucideIcon;
  readonly error?: string;
  readonly outerClassName?: string;
  readonly groupLabel?: boolean;
  readonly labelAction?: ReactNode;
  readonly trailingAction?: ReactNode;
  readonly dataTestId?: string;
}
export type CredentialFieldComponent = (props: CredentialFieldProps) => ReactElement;
