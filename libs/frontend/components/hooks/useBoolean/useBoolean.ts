import { useCallback, useState } from 'react';

import type { UseBoolean } from './types/useBoolean';

/** Stable, reusable boolean state without browser or business dependencies. */
export const useBoolean: UseBoolean = (initial = false) => {
  const [value, setValue] = useState(initial);

  const setTrue = useCallback(() => {
    setValue(true);
  }, []);

  const setFalse = useCallback(() => {
    setValue(false);
  }, []);

  const toggle = useCallback(() => {
    setValue((current) => !current);
  }, []);

  return {
    value,
    setTrue,
    setFalse,
    toggle,
    setValue,
  };
};
