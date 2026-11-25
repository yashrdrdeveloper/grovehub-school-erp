import { Box } from '../Box';
import { Date } from '../Date';

import styles from './DateRange.module.css';
import type { DateRangeProps } from './types/dateRange.props';

export const DateRange = ({ value, onChange, disabled, dataTestId }: DateRangeProps) => (
  <Box className={styles.range} dataTestId={dataTestId}>
    <Date label={__('Start date')} value={value.startDate} disabled={disabled} onChange={(next) => onChange({ ...value, startDate: next })} />
    <Date label={__('End date')} value={value.endDate} disabled={disabled} onChange={(next) => onChange({ ...value, endDate: next })} />
  </Box>
);
