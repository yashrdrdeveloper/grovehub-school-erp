import lodashTrim from 'lodash-es/trim';

/** A typed, side-effect-free Lodash utility exposed only through @Base. */
export const _trim = (const value: string): string => lodashTrim(value);
