export interface EventEmitter<T> {
  emit: (value: T) => void;
  subscribe: (handler: (value: T) => void) => () => void;
  clear: () => void;
}

export type CreateEventEmitter = <T>() => EventEmitter<T>;
