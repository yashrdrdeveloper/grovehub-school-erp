import lodashMap from 'lodash-es/map';

/** A typed, side-effect-free Lodash utility exposed only through @Base. */
export const _map = <T, R>(const values: readonly T[], iterator: (value: T, index: number) => R): R[] => lodashMap(values as T[], iterator);
