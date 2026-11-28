import type { useDateRange } from '../../../hooks/useDateRange';
import type { DateRangeProps } from '../../../types/dateRange.props';

export type DateRangePresetsProps = Pick<DateRangeProps, 'disabled'> & {
  readonly controls: ReturnType<typeof useDateRange>;
};
