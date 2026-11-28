export type ReadProperty = <T, K extends keyof T>(key: K) => (value: T) => T[K];
