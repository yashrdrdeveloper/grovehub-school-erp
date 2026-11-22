import * as DialogPrimitive from '@radix-ui/react-dialog';
import { X } from 'lucide-react';
import { cn } from '@Base/utilities/cn';
import type { DialogComponent } from './types/dialog.props';

export const Dialog: DialogComponent = ({
  open, onOpenChange, title, description, children, className, dataTestId,
}) => {
  let closeTestId: string | undefined;
  if (dataTestId) closeTestId = dataTestId + '-close';
  return (
  <DialogPrimitive.Root open={open} onOpenChange={onOpenChange}>
    <DialogPrimitive.Portal>
      <DialogPrimitive.Overlay className="dialog-overlay"/>
      <DialogPrimitive.Content className={cn('dialog-content rounded-xl border border-border bg-background text-foreground shadow-xl', className)} data-testid={dataTestId}>
        <div className="dialog-heading">
          <div>
            <DialogPrimitive.Title>{title}</DialogPrimitive.Title>
            {description && <DialogPrimitive.Description>{description}</DialogPrimitive.Description>}
          </div>
          <DialogPrimitive.Close className="icon-button" data-testid={closeTestId} aria-label="Close"><X size={19}/></DialogPrimitive.Close>
        </div>
        {children}
      </DialogPrimitive.Content>
    </DialogPrimitive.Portal>
  </DialogPrimitive.Root>
  );
};
