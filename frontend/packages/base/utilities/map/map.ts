import lodashMap from 'lodash-es/map';
import type { MapValues } from './types/map';

export const _map: MapValues = (values, iterator) => lodashMap(values as unknown[], iterator);
