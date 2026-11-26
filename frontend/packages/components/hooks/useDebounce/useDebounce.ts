import { useCallback, useEffect, useRef } from 'react';

import type { UseDebounceResult } from './types/useDebounce';

export const useDebounce = <T>(
  callback: (value: T) => void,
  delay: number,
): UseDebounceResult<T> => {
  const latest = useRef(callback);
  const timer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const pending = useRef<{ value: T } | undefined>(undefined);

  useEffect(() => {
    latest.current = callback;
  }, [callback]);

  const cancel = useCallback(() => {
    clearTimeout(timer.current);
    timer.current = undefined;
    pending.current = undefined;
  }, []);

  const flush = useCallback(() => {
    const entry = pending.current;
    cancel();
    if (entry) latest.current(entry.value);
  }, [cancel]);

  const execute = useCallback(
    (value: T) => {
      cancel();
      pending.current = { value };
      timer.current = setTimeout(flush, Math.max(0, delay));
    },
    [cancel, delay, flush],
  );

  useEffect(() => cancel, [cancel, delay]);

  return { execute, cancel, flush };
};
