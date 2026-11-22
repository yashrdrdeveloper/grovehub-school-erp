import { useSyncExternalStore } from 'react';

import { getMobileSnapshot } from './helpers/loginViewport.getMobileSnapshot';
import { subscribeToViewport } from './helpers/loginViewport.subscribeToViewport';
import type { UseLoginViewport } from './types/useLoginViewport';

/** Only a hook: viewport reading and subscription are named reusable helpers. */
export const useLoginViewport: UseLoginViewport = () =>
  useSyncExternalStore(subscribeToViewport, getMobileSnapshot, () => false);
