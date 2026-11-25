import { useEffect, useRef } from 'react';

import type { TableSelectionControlProps } from './types/tableSelection.props';
import styles from './TableSelection.module.css';

export const TableSelection = ({ checked, indeterminate, disabled, label, onChange }: TableSelectionControlProps) => {
  const ref = useRef<HTMLInputElement>(null);
  useEffect(() => {
    if (ref.current) ref.current.indeterminate = Boolean(indeterminate);
  }, [indeterminate]);
  return (
    <input
      ref={ref}
      className={styles.checkbox}
      type="checkbox"
      checked={checked}
      disabled={disabled}
      aria-label={label}
      onChange={onChange}
      onClick={(event) => event.stopPropagation()}
    />
  );
};
