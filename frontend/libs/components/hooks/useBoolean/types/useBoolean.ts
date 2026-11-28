import type { Dispatch, SetStateAction } from 'react';

export interface BooleanState {
  readonly value: boolean;
  readonly setTrue: () => void;
  readonly setFalse: () => void;
  readonly toggle: () => void;
  readonly setValue: Dispatch<SetStateAction<boolean>>;
}
export type UseBoolean = (initial?: boolean) => BooleanState;
