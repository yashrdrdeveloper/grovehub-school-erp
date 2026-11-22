import { useSyncExternalStore } from 'react';

import { getMobileSnapshot } from './helpers/getMobileSnapshot';
import { subscribeToViewport } from './helpers/subscribeToViewport';
import type { UseLoginViewport } from './types/useLoginViewport';

/** Only a hook: viewport reading and subscription are named reusable helpers. */
export const useLoginViewport: UseLoginViewport = () =>
  useSyncExternalStore(subscribeToViewport, getMobileSnapshot, () => false);
