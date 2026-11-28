import type { ReadProperty } from './types/readProperty';

/** Pure reusable typed selector factory. */
export const readProperty: ReadProperty = (key) => (value) => value[key];
