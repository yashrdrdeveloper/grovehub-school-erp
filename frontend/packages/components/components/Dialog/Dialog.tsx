import { BOX_TAG } from '@Components/components/Box';
import { X } from '@Components/icons';
import { cn } from '@Base/utilities/cn';

import { Box } from '../Box';
import { DialogRoot } from './components/root';
import { DialogPortal } from './components/portal';
import { DialogOverlay } from './components/overlay';
import { DialogContent } from './components/content';
import { DialogTitle } from './components/title';
import { DialogDescription } from './components/description';
import { DialogClose } from './components/close';

import type { DialogComponent } from './types/dialog.props';

/** Accessible shared Dialog composed from seven separately owned Radix primitives. */
export const Dialog: DialogComponent = ({
  open,
  onOpenChange,
  title,
  description,
  children,
  className,
  dataTestId,
}) => {
  const closeTestId = dataTestId && `${dataTestId}-close`;

  return (
    <DialogRoot open={open} onOpenChange={onOpenChange}>
      <DialogPortal>
        <DialogOverlay className="dialog-overlay" />
        <DialogContent
          className={cn(
            'dialog-content rounded-xl border border-border bg-background text-foreground shadow-xl',
            className,
          )}
          data-testid={dataTestId}
        >
          <Box as={BOX_TAG.DIV} className="dialog-heading">
            <Box as={BOX_TAG.DIV}>
              <DialogTitle>{title}</DialogTitle>
              {description && <DialogDescription>{description}</DialogDescription>}
            </Box>
            <DialogClose
              className="icon-button"
              data-testid={closeTestId || undefined}
              aria-label="Close"
            >
              <X size={19} aria-hidden="true" />
            </DialogClose>
          </Box>
          {children}
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  );
};
