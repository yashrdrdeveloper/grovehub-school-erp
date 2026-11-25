import { cn } from '@Base/utilities/cn';

import { Dialog } from '../Dialog';

import styles from './Sheet.module.css';
import type { SheetProps } from './types/sheet.props';

export const Sheet = ({ className, ...props }: SheetProps) => (
  <Dialog {...props} className={cn(styles.sheet, className)} />
);
