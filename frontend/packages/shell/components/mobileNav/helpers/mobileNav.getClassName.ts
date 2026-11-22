import { cn } from '@Base/utilities/cn';

import styles from '../MobileNav.module.css';

import type { GetMobileNavClassName } from './types/mobileNav.getClassName';

/** The router receives a stable class resolver instead of inline CSS logic. */
export const getMobileNavClassName: GetMobileNavClassName = ({ isActive }) =>
  cn(styles.item, isActive && styles.active);
