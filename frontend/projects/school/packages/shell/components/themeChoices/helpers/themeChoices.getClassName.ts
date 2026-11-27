import { cn } from '@Base/utilities/cn';

import styles from '../ThemeChoices.module.css';

import type { GetThemeChoiceClassName } from './types/themeChoices.getClassName';

/** CSS module modifier belongs to this picker; no reusable global CSS class escapes. */
export const getThemeChoiceClassName: GetThemeChoiceClassName = (active) =>
  cn(styles.item, active && styles.active);
