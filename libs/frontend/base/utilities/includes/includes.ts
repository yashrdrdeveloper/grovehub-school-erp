import lodashIncludes from 'lodash-es/includes';
import type { Includes } from './types/includes';

export const _includes: Includes = (values, value) => lodashIncludes(values, value);
