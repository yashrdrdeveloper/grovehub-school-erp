import { useEffect, useRef, useState } from 'react';

import { useDebounce } from '@Components/hooks/useDebounce';

import { DebouncedInputView } from '../../DebouncedInput';
import { DEBOUNCED_INPUT_DEFAULTS } from '../../constants/debouncedInput.defaults';
import type { DebouncedInputContainerComponent } from './types/debouncedInputContainer';

export const DebouncedInputContainer: DebouncedInputContainerComponent = ({
  value,
  onChange,
  onImmediateChange,
  debounceMs = DEBOUNCED_INPUT_DEFAULTS.debounceMs,
  minLength = DEBOUNCED_INPUT_DEFAULTS.minLength,
  flushOnBlur = DEBOUNCED_INPUT_DEFAULTS.flushOnBlur,
  flushOnEnter = DEBOUNCED_INPUT_DEFAULTS.flushOnEnter,
  immediateOnClear = DEBOUNCED_INPUT_DEFAULTS.immediateOnClear,
  onBlur,
  onKeyDown,
  onCompositionStart,
  onCompositionEnd,
  ...inputProps
}) => {
  const [draft, setDraft] = useState(value);
  const composing = useRef(false);
  const lastEmitted = useRef(value);
  const { execute, cancel, flush } = useDebounce((next: string) => {
    if (next === lastEmitted.current) return;
    lastEmitted.current = next;
    onChange(next);
  }, debounceMs);

  useEffect(() => {
    cancel();
    lastEmitted.current = value;
    setDraft(value);
  }, [value, cancel]);

  const schedule = (next: string) => {
    cancel();
    if (composing.current) return;
    if (next === lastEmitted.current) return;
    if (next === '' && immediateOnClear) {
      lastEmitted.current = next;
      onChange(next);
      return;
    }
    if (next.length >= minLength) execute(next);
  };

  return (
    <DebouncedInputView
      {...inputProps}
      value={draft}
      onChange={(event) => {
        const next = event.currentTarget.value;
        setDraft(next);
        onImmediateChange?.(next);
        schedule(next);
      }}
      onCompositionStart={(event) => {
        composing.current = true;
        cancel();
        onCompositionStart?.(event);
      }}
      onCompositionEnd={(event) => {
        composing.current = false;
        schedule(event.currentTarget.value);
        onCompositionEnd?.(event);
      }}
      onBlur={(event) => {
        if (flushOnBlur && !composing.current) flush();
        onBlur?.(event);
      }}
      onKeyDown={(event) => {
        if (event.key === 'Enter' && flushOnEnter && !composing.current) flush();
        onKeyDown?.(event);
      }}
    />
  );
};
