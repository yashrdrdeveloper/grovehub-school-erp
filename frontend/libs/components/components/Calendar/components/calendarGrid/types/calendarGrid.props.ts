import type { CalendarProps } from '../../../types/calendar.props';
import type { CalendarMode } from '../../../types/calendar.mode';
import type { ButtonProps } from '../../../../Button';

export type CalendarGridProps = Pick<
  CalendarProps,
  'value' | 'range' | 'disabled' | 'onChange' | 'onDayHover'
> & {
  readonly mode: CalendarMode;
  readonly month: Date;
  readonly size: ButtonProps['size'];
  readonly unavailable: (date: string) => boolean;
  readonly chooseMonth: (index: number) => void;
  readonly chooseYear: (year: number) => void;
};
