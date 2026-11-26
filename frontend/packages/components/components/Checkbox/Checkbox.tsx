import type { CheckboxProps } from './types/checkbox.props';
import styles from './Checkbox.module.css';

export const Checkbox = ({ checked, onChange, label, disabled, dataTestId }: CheckboxProps) => (
  <input
    type="checkbox"
    className={styles.checkbox}
    checked={checked}
    disabled={disabled}
    aria-label={label}
    data-testid={dataTestId}
    onChange={(event) => onChange(event.currentTarget.checked)}
  />
);
