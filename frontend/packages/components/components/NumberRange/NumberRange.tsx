import { Box } from '../Box';
import { Number } from '../Number';

import styles from './NumberRange.module.css';
import type { NumberRangeProps } from './types/numberRange.props';

export const NumberRange = ({ value, onChange, disabled, dataTestId }: NumberRangeProps) => (
  <Box className={styles.range} dataTestId={dataTestId}>
    <Number label={__('Minimum')} value={value.min} disabled={disabled} onChange={(next) => onChange({ ...value, min: next })} />
    <Number label={__('Maximum')} value={value.max} disabled={disabled} onChange={(next) => onChange({ ...value, max: next })} />
  </Box>
);
