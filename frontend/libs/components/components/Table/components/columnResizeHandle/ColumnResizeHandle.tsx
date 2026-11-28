import { useRef } from 'react';
import { Box } from '../../../Box';
import type { ColumnResizeHandleProps } from './types/columnResizeHandle.props';
import styles from './ColumnResizeHandle.module.css';

export const ColumnResizeHandle = ({
  label,
  width,
  min,
  max,
  onResize,
}: ColumnResizeHandleProps) => {
  const drag = useRef<{ x: number; width: number; pointerId: number } | null>(null);
  const clamp = (next: number) => Math.max(min, Math.min(max, Math.round(next)));
  return (
    <Box
      className={styles.handle}
      role="separator"
      tabIndex={0}
      aria-label={label}
      aria-orientation="vertical"
      aria-valuemin={min}
      aria-valuemax={max}
      aria-valuenow={width}
      onPointerDown={(event) => {
        if (event.button !== 0) return;
        event.preventDefault();
        event.stopPropagation();
        drag.current = { x: event.clientX, width, pointerId: event.pointerId };
        event.currentTarget.setPointerCapture(event.pointerId);
        event.currentTarget.dataset.resizing = 'true';
      }}
      onPointerMove={(event) => {
        const current = drag.current;
        if (!current || current.pointerId !== event.pointerId) return;
        onResize(clamp(current.width + event.clientX - current.x));
      }}
      onLostPointerCapture={(event) => {
        drag.current = null;
        delete event.currentTarget.dataset.resizing;
      }}
      onKeyDown={(event) => {
        if (event.key === 'ArrowLeft' || event.key === 'ArrowRight') {
          event.preventDefault();
          onResize(clamp(width + (Number(event.key === 'ArrowRight') * 2 - 1) * 10));
        }
      }}
    />
  );
};
