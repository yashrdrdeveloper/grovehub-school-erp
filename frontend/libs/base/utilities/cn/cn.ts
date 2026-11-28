import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

import type { Cn } from './types/cn';

export const cn: Cn = (...values) => twMerge(clsx(values));
