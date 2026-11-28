export type Find = <T>(values: readonly T[], predicate: (value: T) => boolean) => T | undefined;
