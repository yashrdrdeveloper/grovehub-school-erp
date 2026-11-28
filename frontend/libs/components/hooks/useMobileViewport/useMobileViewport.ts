import { useSyncExternalStore } from 'react';

const query = '(max-width: 640px)';
const subscribe = (notify: () => void) => {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function')
    return () => undefined;
  const media = window.matchMedia(query);
  media.addEventListener('change', notify);
  return () => media.removeEventListener('change', notify);
};
const snapshot = () =>
  typeof window !== 'undefined' &&
  typeof window.matchMedia === 'function' &&
  window.matchMedia(query).matches;

export const useMobileViewport = (): boolean =>
  useSyncExternalStore(subscribe, snapshot, () => false);
