import { cn } from '@Base/utilities/cn';

import styles from '../PinnedNavigation.module.css';

import type { GetPinnedAppClassName } from './types/pinnedNavigation.getAppClassName';

/** A single renderer-owned active-state mapper for pinned route links. */
export const getPinnedAppClassName: GetPinnedAppClassName = ({ isActive }) =>
  cn(styles.item, isActive && styles.active);
