import { STYLE_MEDIA_QUERY } from '@Styles/tokens/breakpoints';
import type { SubscribeToViewport } from './types/loginViewport.subscribeToViewport';

export const subscribeToViewport: SubscribeToViewport = notify => {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return () => {};
  const query = window.matchMedia(STYLE_MEDIA_QUERY.PHONE);
  query.addEventListener('change', notify);
  return () => query.removeEventListener('change', notify);
};
