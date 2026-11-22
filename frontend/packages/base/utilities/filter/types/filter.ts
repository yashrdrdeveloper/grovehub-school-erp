export type Filter = <T>(values: readonly T[], predicate: (item: T) => boolean) => T[];
