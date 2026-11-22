import lodashKeys from 'lodash-es/keys';
import type { Keys } from './types/keys';

export const _keys: Keys = value => lodashKeys(value);
