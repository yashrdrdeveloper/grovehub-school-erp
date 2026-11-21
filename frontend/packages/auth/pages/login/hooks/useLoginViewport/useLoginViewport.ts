import { useSyncExternalStore } from 'react';

const PHONE_VIEWPORT = '(max-width: 720px)';

const getMobileSnapshot = (): boolean =>
  typeof window !== 'undefined'
  && typeof window.matchMedia === 'function'
  && window.matchMedia(PHONE_VIEWPORT).matches;

const subscribeToViewport = (notify: () => void): (() => void) => {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return () => {};
  const query = window.matchMedia(PHONE_VIEWPORT);
  query.addEventListener('change', notify);
  return () => query.removeEventListener('change', notify);
};

/** Keep a single active form in the DOM; CSS-only duplication breaks IDs and autofill. */
export const useLoginViewport = (): boolean =>
  useSyncExternalStore(subscribeToViewport, getMobileSnapshot, () => false);
