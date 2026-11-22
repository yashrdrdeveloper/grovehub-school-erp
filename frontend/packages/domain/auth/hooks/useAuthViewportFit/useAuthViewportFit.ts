import { useLayoutEffect, useRef } from 'react';
import { STYLE_MEDIA_QUERY } from '@Styles/tokens/breakpoints';

/**
 * Fit a shared auth state into the usable desktop/tablet panel without clipping.
 *
 * We measure the *actual* natural form height (including errors and recovery
 * states) instead of guessing from CSS breakpoints. Scaling is limited to
 * a readable range; small/zoomed/keyboard-constrained views retain scrolling.
 * Mobile (< 721px) uses the dedicated compact layout and is never scaled.
 */
export const useAuthViewportFit = () => {
  const contentRef = useRef<HTMLElement>(null);
  const toolbarRef = useRef<HTMLDivElement>(null);
  const mainRef = useRef<HTMLDivElement>(null);

  useLayoutEffect(() => {
    const content = contentRef.current;
    const toolbar = toolbarRef.current;
    const main = mainRef.current;
    if (!content || !toolbar || !main) return;

    let desktop: MediaQueryList | null = null;
    if (typeof window.matchMedia === 'function') desktop = window.matchMedia(STYLE_MEDIA_QUERY.DESKTOP);
    const MIN_READABLE_SCALE = 0.82;
    let frame = 0;
    let active = true;

    const measure = () => {
      frame = 0;
      if (!desktop?.matches) {
        content.removeAttribute('data-auth-fit');
        content.style.removeProperty('--auth-panel-scale');
        content.style.removeProperty('--auth-panel-center');
        return;
      }

      const style = window.getComputedStyle(content);
      const insetBottom = Number.parseFloat(style.paddingBottom) || 0;
      const insetTop = Number.parseFloat(style.paddingTop) || 0;
      const top = Math.max(insetTop, toolbar.offsetTop + toolbar.offsetHeight + 12);
      const bottom = Math.max(top, content.clientHeight - insetBottom - 10);
      const space = Math.max(0, bottom - top);
      // offsetHeight ignores CSS transforms, so the calculation is stable even
      // when the visual content has already been scaled by a previous measure.
      const naturalHeight = Math.max(1, main.offsetHeight);
      const scale = Math.min(1, space / (naturalHeight + 4));

      if (scale < MIN_READABLE_SCALE || space < 240) {
        content.setAttribute('data-auth-fit', 'scroll');
        content.style.removeProperty('--auth-panel-scale');
        content.style.removeProperty('--auth-panel-center');
        return;
      }

      content.style.setProperty('--auth-panel-scale', String(Math.max(MIN_READABLE_SCALE, scale)));
      content.style.setProperty('--auth-panel-center', `${(top + bottom) / 2}px`);
      content.setAttribute('data-auth-fit', 'fit');
    };

    const schedule = () => {
      if (!active) return;
      if (frame) window.cancelAnimationFrame(frame);
      frame = window.requestAnimationFrame(measure);
    };

    measure();
    let observer: ResizeObserver | null = null;
    if (typeof ResizeObserver !== 'undefined') observer = new ResizeObserver(schedule);
    observer?.observe(content);
    observer?.observe(toolbar);
    observer?.observe(main);
    window.addEventListener('resize', schedule);
    desktop?.addEventListener?.('change', schedule);

    // Web fonts can change form height after the first layout.
    if (document.fonts?.ready) void document.fonts.ready.then(schedule);

    return () => {
      active = false;
      if (frame) window.cancelAnimationFrame(frame);
      observer?.disconnect();
      window.removeEventListener('resize', schedule);
      desktop?.removeEventListener?.('change', schedule);
    };
  }, []);

  return { contentRef, toolbarRef, mainRef };
};
