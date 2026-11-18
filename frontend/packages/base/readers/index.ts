export const readProperty = <T, K extends keyof T>(key: K) =>
  (value: T): T[K] => value[key];
