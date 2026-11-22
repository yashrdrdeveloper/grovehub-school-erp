export interface BooleanState {
  readonly value: boolean;
  readonly setTrue: () => void;
  readonly setFalse: () => void;
  readonly toggle: () => void;
  readonly setValue: (next: boolean) => void;
}
export type UseBoolean = (initial?: boolean) => BooleanState;
