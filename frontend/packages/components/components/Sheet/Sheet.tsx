import { cn } from '@Base/utilities/cn';

import { Dialog } from '../Dialog';

import styles from './Sheet.module.css';
import type { SheetProps } from './types/sheet.props';

export const Sheet = ({ className, dataTestId, ...props }: SheetProps) => (
  <Dialog {...props} dataTestId={dataTestId} className={cn(styles.sheet, className)} />
);
