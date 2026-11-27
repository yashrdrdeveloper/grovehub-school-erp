import { useSyncExternalStore } from 'react';

const query = '(max-width: 760px)';
function subscribe(notify: () => void) {
  const media = window.matchMedia(query);
  media.addEventListener('change', notify);
  return () => media.removeEventListener('change', notify);
}
function snapshot() { return window.matchMedia(query).matches; }
function serverSnapshot() { return false; }

export const useLoginViewport = () => useSyncExternalStore(subscribe, snapshot, serverSnapshot);
