import { useLayoutEffect, useRef, useState } from 'react';

import { Box } from '../../../Box';

import { getScrollbarGeometry, getScrollbarOffset } from './helpers/tableScrollbar.geometry';
import type { TableScrollbarProps } from './types/tableScrollbar.props';
import styles from './TableScrollbar.module.css';

export const TableScrollbar = ({ viewport, viewportId, axis }: TableScrollbarProps) => {
  const vertical = axis === 'vertical';
  const [geometry, setGeometry] = useState(() => getScrollbarGeometry(0, 0, 0));
  const [offset, setOffset] = useState(0);
  const [dragging, setDragging] = useState(false);
  const drag = useRef<{ pointerId: number; grab: number } | null>(null);

  useLayoutEffect(() => {
    const element = viewport.current;
    if (!element) return;
    const measure = () => {
      const nextOffset = (vertical && element.scrollTop) || (!vertical && element.scrollLeft) || 0;
      const next = getScrollbarGeometry(
        (vertical && element.clientHeight) || (!vertical && element.clientWidth) || 0,
        (vertical && element.scrollHeight) || (!vertical && element.scrollWidth) || 0,
        nextOffset,
      );
      setGeometry((current) => {
        if (current.maximum === next.maximum && current.size === next.size && current.position === next.position) return current;
        return next;
      });
      setOffset(nextOffset);
    };
    measure();
    element.addEventListener('scroll', measure, { passive: true });
    const observer = typeof ResizeObserver !== 'undefined' && new ResizeObserver(measure);
    if (observer) {
      observer.observe(element);
      const content = element.firstElementChild;
      if (content) observer.observe(content);
    }
    window.addEventListener('resize', measure);
    return () => {
      element.removeEventListener('scroll', measure);
      window.removeEventListener('resize', measure);
      if (observer) observer.disconnect();
    };
  }, [viewport, vertical]);

  const scrollTo = (next: number) => {
    const element = viewport.current;
    if (!element) return;
    const position = Math.min(geometry.maximum, Math.max(0, next));
    if (vertical) element.scrollTop = position;
    else element.scrollLeft = position;
  };

  if (!geometry.maximum) return null;

  const getThumbStyle = () => {
    if (vertical) return { height: geometry.size + '%', top: geometry.position + '%' };
    return { width: geometry.size + '%', left: geometry.position + '%' };
  };
  const getLabel = () => {
    if (vertical) return __('Vertical table scrollbar');
    return __('Horizontal table scrollbar');
  };
  const getPointer = (rect: DOMRect, x: number, y: number) => {
    if (vertical) return { point: y - rect.top, length: rect.height };
    return { point: x - rect.left, length: rect.width };
  };
  const getPage = () => {
    if (vertical) return viewport.current?.clientHeight ?? 0;
    return viewport.current?.clientWidth ?? 0;
  };
  const getPreviousKey = () => {
    if (vertical) return 'ArrowUp';
    return 'ArrowLeft';
  };
  const getNextKey = () => {
    if (vertical) return 'ArrowDown';
    return 'ArrowRight';
  };

  return (
    <Box
      className={styles.track}
      data-table-scrollbar={axis}
      data-dragging={dragging}
      role="scrollbar"
      aria-label={getLabel()}
      aria-controls={viewportId}
      aria-orientation={axis}
      aria-valuemin={0}
      aria-valuemax={Math.round(geometry.maximum)}
      aria-valuenow={Math.round(Math.min(geometry.maximum, Math.max(0, offset)))}
      tabIndex={0}
      onPointerDown={(event) => {
        if (event.button !== 0) return;
        const rect = event.currentTarget.getBoundingClientRect();
        const { point, length } = getPointer(rect, event.clientX, event.clientY);
        const thumbStart = geometry.position / 100 * length;
        const thumbLength = geometry.size / 100 * length;
        const inside = point >= thumbStart && point <= thumbStart + thumbLength;
        const grab = (inside && point - thumbStart) || (!inside && thumbLength / 2) || 0;
        event.preventDefault();
        event.currentTarget.focus({ preventScroll: true });
        event.currentTarget.setPointerCapture(event.pointerId);
        drag.current = { pointerId: event.pointerId, grab };
        setDragging(true);
        if (!inside) scrollTo(getScrollbarOffset(point, length, geometry.size, grab, geometry.maximum));
      }}
      onPointerMove={(event) => {
        if (!drag.current || drag.current.pointerId !== event.pointerId) return;
        const rect = event.currentTarget.getBoundingClientRect();
        const { point, length } = getPointer(rect, event.clientX, event.clientY);
        scrollTo(getScrollbarOffset(point, length, geometry.size, drag.current.grab, geometry.maximum));
      }}
      onLostPointerCapture={() => { drag.current = null; setDragging(false); }}
      onKeyDown={(event) => {
        const page = getPage();
        if (event.key === 'Home') scrollTo(0);
        else if (event.key === 'End') scrollTo(geometry.maximum);
        else if (event.key === 'PageUp') scrollTo(offset - page);
        else if (event.key === 'PageDown') scrollTo(offset + page);
        else if (event.key === getPreviousKey()) scrollTo(offset - 40);
        else if (event.key === getNextKey()) scrollTo(offset + 40);
        else return;
        event.preventDefault();
      }}
    >
      <Box className={styles.thumb} style={getThumbStyle()} />
    </Box>
  );
};
