import { useEffect, useState } from 'react';

import { __ as resolveMessage } from '@Base/i18n';

import type { UseInitI18n } from './types/useInitI18n';

/**
 * Single browser installation boundary for the app-wide message resolver.
 * Routes must not mount until the effect has assigned window.__.
 */
export const useInitI18n: UseInitI18n = () => {
  const [ready, setReady] = useState(false);

  useEffect(() => {
    if (typeof window === 'undefined') return;

    // Repeated assignment is safe under React StrictMode effect replay.
    window.__ = resolveMessage;
    setReady(true);
  }, []);

  return ready;
};
