import { useEffect, useState } from 'react';

import { installShellI18n } from '../../helpers/shellContainer.installI18n';

import { InitI18nStatus } from './constants/useInitI18n.status';
import type { UseInitI18n } from './types/useInitI18n';

/** Single, StrictMode-safe browser installer for the frontend-wide global __. */
export const useInitI18n: UseInitI18n = () => {
  const [status, setStatus] = useState<InitI18nStatus>(InitI18nStatus.INITIALIZING);

  useEffect(() => {
    if (typeof window === 'undefined') return;

    try {
      installShellI18n();
      setStatus(InitI18nStatus.READY);
    } catch {
      setStatus(InitI18nStatus.FAILED);
    }
  }, []);

  return status;
};
