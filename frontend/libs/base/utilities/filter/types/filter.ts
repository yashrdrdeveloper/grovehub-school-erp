export interface Filter {
  <T, S extends T>(values: readonly T[], predicate: (item: T) => item is S): S[];
  <T>(values: readonly T[], predicate: (item: T) => boolean): T[];
}
