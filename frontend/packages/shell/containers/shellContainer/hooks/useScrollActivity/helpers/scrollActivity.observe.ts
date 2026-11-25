const SCROLL_IDLE_DELAY_MS = 1000;
const SCROLL_ACTIVE_ATTRIBUTE = 'data-gh-scroll-active';

/** Capture native scroll events from the document and all nested scroll owners. */
export const observeScrollActivity = (document: Document): (() => void) => {
  const pending = new Map<Element, ReturnType<typeof setTimeout>>();

  const clear = (element: Element) => {
    const timeout = pending.get(element);
    if (timeout !== undefined) clearTimeout(timeout);
    pending.delete(element);
    element.removeAttribute(SCROLL_ACTIVE_ATTRIBUTE);
  };

  const onScroll = (event: Event) => {
    const target = event.target;
    const element = target === document && document.documentElement || target;
    if (!(element instanceof Element)) return;
    const timeout = pending.get(element);
    if (timeout !== undefined) clearTimeout(timeout);
    element.setAttribute(SCROLL_ACTIVE_ATTRIBUTE, 'true');
    pending.set(element, setTimeout(() => clear(element), SCROLL_IDLE_DELAY_MS));
  };

  document.addEventListener('scroll', onScroll, { capture: true, passive: true });

  return () => {
    document.removeEventListener('scroll', onScroll, true);
    pending.forEach((_timeout, element) => clear(element));
  };
};
