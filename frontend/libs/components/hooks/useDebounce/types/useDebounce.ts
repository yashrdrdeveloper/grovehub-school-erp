export interface UseDebounceResult<T> {
  readonly execute: (value: T) => void;
  readonly cancel: () => void;
  readonly flush: () => void;
}
