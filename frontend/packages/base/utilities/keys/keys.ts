import lodashKeys from 'lodash-es/keys';

/** A typed, side-effect-free Lodash utility exposed only through @Base. */
export const _keys = (const value: object): string[] => lodashKeys(value);
