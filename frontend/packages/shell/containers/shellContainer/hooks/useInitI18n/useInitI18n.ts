import { useEffect, useState } from 'react';

import { __ as resolveMessage } from '@Base/i18n';

import { InitI18nStatus } from './constants/useInitI18n.status';
import type { UseInitI18n } from './types/useInitI18n';

/** Single, StrictMode-safe browser installer for the frontend-wide global __. */
export const useInitI18n: UseInitI18n = () => {
  const [status, setStatus] = useState(InitI18nStatus.INITIALIZING);

  useEffect(() => {
    if (typeof window === 'undefined') return;

    try {
      window.__ = resolveMessage;
      if (window.__ !== resolveMessage) {
        throw new Error('Failed to register the global message formatter');
      }
      setStatus(InitI18nStatus.READY);
    } catch {
      setStatus(InitI18nStatus.FAILED);
    }
  }, []);

  return status;
};
