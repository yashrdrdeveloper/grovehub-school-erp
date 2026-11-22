import { PHONE_VIEWPORT } from '../../../../config';
import type { SubscribeToViewport } from './types/subscribeToViewport';

export const subscribeToViewport: SubscribeToViewport = notify => {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return () => {};
  const query = window.matchMedia(PHONE_VIEWPORT);
  query.addEventListener('change', notify);
  return () => query.removeEventListener('change', notify);
};
