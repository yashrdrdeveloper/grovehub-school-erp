import type { CreateEventEmitter } from './types/createEventEmitter';

/** Technical event emitter with no business/domain dependencies. */
export const createEventEmitter: CreateEventEmitter = <T>() => {
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
