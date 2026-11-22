import lodashFind from 'lodash-es/find';

/** A typed, side-effect-free Lodash utility exposed only through @Base. */
export const _find = <T>(const values: readonly T[], predicate: (value: T) => boolean): T | undefined => lodashFind(values as T[], predicate);
