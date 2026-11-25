import { ButtonWithIcon } from '../ButtonWithIcon';

import type { IconButtonProps } from './types/iconButton.props';
import styles from './IconButton.module.css';

export const IconButton = ({ label, dataTestId, ...props }: IconButtonProps) => (
  <ButtonWithIcon {...props} dataTestId={dataTestId} aria-label={label} title={label} className={styles.button} />
);
