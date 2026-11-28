import type {
  AnchorHTMLAttributes,
  ForwardRefExoticComponent,
  PropsWithoutRef,
  RefAttributes,
} from 'react';
export interface ExternalAnchorProps extends AnchorHTMLAttributes<HTMLAnchorElement> {
  readonly dataTestId?: string;
}
export type ExternalAnchorComponent = ForwardRefExoticComponent<
  PropsWithoutRef<ExternalAnchorProps> & RefAttributes<HTMLAnchorElement>
>;
