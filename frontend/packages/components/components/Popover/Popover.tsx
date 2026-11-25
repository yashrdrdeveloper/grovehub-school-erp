import { useLayoutEffect, useRef, useState } from 'react';
import { X } from '@Components/icons';

import { Box } from '../Box';
import { DialogRoot, DialogPortal, DialogContent, DialogTitle, DialogClose } from '../Dialog';

import { getPopoverPosition } from './helpers/popover.getPosition';
import styles from './Popover.module.css';
import type { PopoverProps } from './types/popover.props';

export const Popover = ({ open, onOpenChange, title, trigger, children, width = 360, fullWidth, compact = false, desktopScrollLayout = false, desktopDropdown = false, desktopDropdownWidth, dataTestId }: PopoverProps) => {
  const anchor = useRef<HTMLElement>(null);
  const [panel, setPanel] = useState<HTMLDivElement | null>(null);
  const [position, setPosition] = useState({ left: 0, top: 0, width, maxHeight: 0 });
  useLayoutEffect(() => {
    if (!open || !panel) return;
    const place = () => {
      const rect = anchor.current?.getBoundingClientRect();
      if (!rect) return;
      const available = Math.max(0, window.innerWidth - 24);
      const getWidth = () => {
        if (desktopDropdown && window.innerWidth > 640) return Math.max(desktopDropdownWidth ?? 88, rect.width);
        return Math.max(width, rect.width);
      };
      const panelWidth = Math.min(available, getWidth());
      panel.style.width = String(panelWidth) + 'px';
      const panelHeight = panel.offsetHeight;
      setPosition(getPopoverPosition({ anchor: rect, panelHeight, viewportWidth: window.innerWidth, viewportHeight: window.innerHeight, width: panelWidth }));
    };
    place();
    const observer = typeof ResizeObserver !== 'undefined' && new ResizeObserver(place);
    if (observer) { observer.observe(panel); if (anchor.current) observer.observe(anchor.current); }
    window.addEventListener('resize', place);
    window.addEventListener('scroll', place, true);
    return () => {
      if (observer) observer.disconnect();
      window.removeEventListener('resize', place);
      window.removeEventListener('scroll', place, true);
    };
  }, [open, panel, width, desktopDropdown, desktopDropdownWidth]);
  return (
    <DialogRoot open={open} onOpenChange={onOpenChange}>
      <Box ref={anchor} dataTestId={dataTestId} className={styles.anchor} data-full-width={fullWidth}>
        {trigger}
      </Box>
      <DialogPortal>
        <DialogContent
          ref={setPanel}
          aria-describedby={undefined}
          className={styles.panel}
          data-compact={compact}
          data-desktop-scroll-layout={desktopScrollLayout}
          data-desktop-dropdown={desktopDropdown}
          style={{ left: position.left, top: position.top, width: position.width }}
          onOpenAutoFocus={(event) => {
            if (!desktopDropdown || window.innerWidth <= 640) return;
            event.preventDefault();
            if (event.target instanceof HTMLElement) event.target.focus({ preventScroll: true });
          }}
          tabIndex={-1}
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
