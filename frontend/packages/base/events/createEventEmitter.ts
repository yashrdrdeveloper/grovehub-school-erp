// Pure technical event emitter. No application or domain dependencies.
export interface EventEmitter<T> {
  emit: (value: T) => void;
  subscribe: (handler: (value: T) => void) => () => void;
  clear: () => void;
}

export const createEventEmitter = <T>(): EventEmitter<T> => {
  const handlers = new Set<(value: T) => void>();

  return {
    emit: (value) => {
      for (const handler of handlers) handler(value);
    },
    subscribe: (handler) => {
      handlers.add(handler);
      return () => { handlers.delete(handler); };
    },
    clear: () => handlers.clear(),
  };
};
