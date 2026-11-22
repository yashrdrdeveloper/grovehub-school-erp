import type { InputHTMLAttributes, ReactElement, ReactNode } from 'react';
import type { LucideIcon } from 'lucide-react';

export interface CredentialFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  readonly id: string;
  readonly label: string;
  readonly Icon: LucideIcon;
  readonly error?: string;
  readonly outerClassName?: string;
  readonly groupLabel?: boolean;
  readonly labelAction?: ReactNode;
  readonly trailingAction?: ReactNode;
}
export type CredentialFieldComponent = (props: CredentialFieldProps) => ReactElement;
