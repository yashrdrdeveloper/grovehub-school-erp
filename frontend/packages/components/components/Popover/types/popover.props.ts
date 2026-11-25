import type { ReactNode } from 'react';

export interface PopoverProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  /** Optional override for primary pointer dismissal; other close actions retain their existing behavior. */
  readonly onOutsideClick?: () => void;
  readonly title: string;
  /** Custom visible heading content; title remains the accessible dialog name. */
  readonly headerContent?: ReactNode;
  /** Hide the visual heading and close button, retaining the accessible dialog name. */
  readonly showHeader?: boolean;
  readonly trigger: ReactNode;
  readonly children: ReactNode;
  readonly width?: number;
  readonly fullWidth?: boolean;
  readonly compact?: boolean;
  /** Desktop-only bounded layout for content with its own scrollable body and footer. */
  readonly desktopScrollLayout?: boolean;
  /** Anchor-sized desktop selection menu; retain the titled mobile sheet. */
  readonly desktopDropdown?: boolean;
  readonly desktopDropdownWidth?: number;
  readonly dataTestId?: string;
}
