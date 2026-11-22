import * as DialogPrimitive from '@radix-ui/react-dialog';
import { X } from 'lucide-react';
import { cn } from '@Base/utilities';
import type { DialogComponent } from './types/dialog.props';

export const Dialog: DialogComponent = ({
  open, onOpenChange, title, description, children, className,
}) => (
  <DialogPrimitive.Root open={open} onOpenChange={onOpenChange}>
    <DialogPrimitive.Portal>
      <DialogPrimitive.Overlay className="dialog-overlay"/>
      <DialogPrimitive.Content className={cn('dialog-content', className)}
        aria-describedby={description ? undefined : undefined}>
        <div className="dialog-heading">
          <div>
            <DialogPrimitive.Title>{title}</DialogPrimitive.Title>
            {description ? <DialogPrimitive.Description>{description}</DialogPrimitive.Description> : null}
          </div>
          <DialogPrimitive.Close className="icon-button" aria-label="Close"><X size={19}/></DialogPrimitive.Close>
        </div>
        {children}
      </DialogPrimitive.Content>
    </DialogPrimitive.Portal>
  </DialogPrimitive.Root>
);
