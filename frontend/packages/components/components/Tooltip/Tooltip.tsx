import { useEffect, useId, useLayoutEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';

import { Box } from '../Box';
import type { TooltipProps } from './types/tooltip.props';
import styles from './Tooltip.module.css';

export const Tooltip = ({ label, children, focusable, fullWidth, dataTestId }: TooltipProps) => {
  const id = useId();
  const bubble = useRef<HTMLElement>(null);
  const anchor = useRef<HTMLElement>(null);
  const [hovered, setHovered] = useState(false);
  const [position, setPosition] = useState<{ left: number; top: number } | null>(null);
  const show = (element: HTMLElement) => {
    if (!label) return;
    const rect = element.getBoundingClientRect();
    setPosition({ left: rect.left + rect.width / 2, top: rect.top - 6 });
  };
  const hide = () => setPosition(null);
  const getTabIndex = () => {
    if (focusable) return 0;
    return undefined;
  };
  useLayoutEffect(() => {
    if (!position || !bubble.current) return;
    const halfWidth = bubble.current.getBoundingClientRect().width / 2;
    const left = Math.min(
      window.innerWidth - halfWidth - 8,
      Math.max(halfWidth + 8, position.left),
    );
    if (left !== position.left) setPosition({ ...position, left });
  }, [position]);
  useEffect(() => {
    if (!position) return;
    const reposition = () => {
      if (anchor.current) show(anchor.current);
    };
    window.addEventListener('scroll', reposition, true);
    window.addEventListener('resize', reposition);
    return () => {
      window.removeEventListener('scroll', reposition, true);
      window.removeEventListener('resize', reposition);
    };
  }, [position]);
  return (
    <Box
      dataTestId={dataTestId}
      ref={anchor}
      className={styles.anchor}
      data-tooltip-hover={hovered}
      data-full-width={fullWidth}
      tabIndex={getTabIndex()}
      aria-label={(focusable && label) || undefined}
      onMouseEnter={(event) => {
        setHovered(true);
        show(event.currentTarget);
      }}
      onMouseLeave={() => {
        setHovered(false);
        hide();
      }}
      onFocus={(event) => {
        if (event.target.matches(':focus-visible')) show(event.currentTarget);
      }}
      onBlur={hide}
      onPointerDown={hide}
      onClick={(event) => {
        if (focusable) show(event.currentTarget);
      }}
      onKeyDown={(event) => {
        if (event.key === 'Escape') hide();
      }}
    >
      {children}
      {label &&
        position &&
        createPortal(
          <Box ref={bubble} id={id} role="tooltip" className={styles.bubble} style={position}>
            {label}
          </Box>,
          document.body,
        )}
    </Box>
  );
};
