import { useEffect } from 'react';

import { observeScrollActivity } from './helpers/scrollActivity.observe';

export const useScrollActivity = (): void => {
  useEffect(() => {
    if (typeof document === 'undefined') return;
    return observeScrollActivity(document);
  }, []);
};
