import { useEffect, useRef, useState } from 'react';
import { X } from '@Components/icons';

import { Box } from '../Box';
import { DialogRoot, DialogPortal, DialogContent, DialogTitle, DialogClose } from '../Dialog';

import styles from './Popover.module.css';
import type { PopoverProps } from './types/popover.props';

export const Popover = ({ open, onOpenChange, title, trigger, children, width = 360, fullWidth, dataTestId }: PopoverProps) => {
  const anchor = useRef<HTMLElement>(null);
  const panel = useRef<HTMLDivElement>(null);
  const [position, setPosition] = useState({ left: 0, top: 0, width });
  useEffect(() => {
    if (!open) return;
    const place = () => {
      const rect = anchor.current?.getBoundingClientRect();
      if (!rect) return;
      const available = Math.max(0, window.innerWidth - 24);
      const panelWidth = Math.min(available, Math.max(width, rect.width));
      const panelHeight = Math.min(panel.current?.offsetHeight ?? 320, window.innerHeight - 24);
      const left = Math.max(12, Math.min(rect.left, window.innerWidth - panelWidth - 12));
      const bottomSpace = window.innerHeight - rect.bottom - 12;
      const top = Math.max(12, Math.min(rect.bottom + 8, window.innerHeight - panelHeight - 12));
      if (bottomSpace < panelHeight && rect.top > bottomSpace) {
        setPosition({ left, top: Math.max(12, rect.top - panelHeight - 8), width: panelWidth });
        return;
      }
      setPosition({ left, top, width: panelWidth });
    };
    place();
    const observer = typeof ResizeObserver !== 'undefined' && new ResizeObserver(place);
    if (observer && panel.current) observer.observe(panel.current);
    window.addEventListener('resize', place);
    window.addEventListener('scroll', place, true);
    return () => {
      if (observer) observer.disconnect();
      window.removeEventListener('resize', place);
      window.removeEventListener('scroll', place, true);
    };
  }, [open, width]);
  return (
    <DialogRoot open={open} onOpenChange={onOpenChange}>
      <Box ref={anchor} dataTestId={dataTestId} className={styles.anchor} data-full-width={fullWidth}>
        {trigger}
      </Box>
      <DialogPortal>
        <DialogContent
          ref={panel}
          aria-describedby={undefined}
          className={styles.panel}
          style={position}
          onCloseAutoFocus={(event) => {
            event.preventDefault();
            anchor.current?.querySelector('button')?.focus();
          }}
        >
          <Box className={styles.heading}>
            <DialogTitle className={styles.title}>{title}</DialogTitle>
            <DialogClose asChild>
              <button type="button" className={styles.close} aria-label={__('Close')}><X size={18} aria-hidden="true" /></button>
            </DialogClose>
          </Box>
          <Box className={styles.content}>{children}</Box>
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  );
};
