import { useCallback, useState } from 'react';
import type { UseBoolean } from './types/useBoolean';

/** Shared boolean state with stable callbacks and no side effects. */
export const useBoolean: UseBoolean = (initial = false) => {
  const [value, setValue] = useState(initial);
  const setTrue = useCallback(() => setValue(true), []);
  const setFalse = useCallback(() => setValue(false), []);
  const toggle = useCallback(() => setValue(current => !current), []);
  return { value, setTrue, setFalse, toggle, setValue };
};
