import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { ArrowLeft, ArrowRight } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';

import { formatCalendarDate, getCalendarDays, parseCalendarDate } from './helpers/calendar.date';
import styles from './Calendar.module.css';
import type { CalendarProps } from './types/calendar.props';

export const Calendar = ({ value, onChange, min, max, disabled, clearable = true, dataTestId }: CalendarProps) => {
  const [month, setMonth] = useState(() => parseCalendarDate(value) ?? new Date());
  const title = new Intl.DateTimeFormat(undefined, { month: 'long', year: 'numeric' }).format(month);
  const today = formatCalendarDate(new Date());
  const days = getCalendarDays(month);
  const shift = (delta: number) => setMonth(new Date(month.getFullYear(), month.getMonth() + delta, 1));
  const unavailable = (date: string) => Boolean(disabled || (min && date < min) || (max && date > max));
  const dayNames = Array.from({ length: 7 }, (_, index) => new Intl.DateTimeFormat(undefined, { weekday: 'short' }).format(new Date(2024, 0, index + 1)));
  return (
    <Box className={styles.calendar} dataTestId={dataTestId}>
      <Box className={styles.heading}>
        <ButtonWithIcon icon={ArrowLeft} variant={BUTTON_VARIANT.GHOST} aria-label={__('Previous month')} onClick={() => shift(-1)} />
        <Box aria-live="polite">{title}</Box>
        <ButtonWithIcon icon={ArrowRight} variant={BUTTON_VARIANT.GHOST} aria-label={__('Next month')} onClick={() => shift(1)} />
      </Box>
      <Box className={styles.days} role="group" aria-label={title}>
        {_map(dayNames, (name, index) => <Box key={index} className={styles.weekday}>{name}</Box>)}
        {_map(days, (day) => {
          const date = formatCalendarDate(day);
          const variant = () => {
            if (date === value) return BUTTON_VARIANT.PRIMARY;
            return BUTTON_VARIANT.GHOST;
          };
          return <Button key={date} variant={variant()} className={styles.day} data-outside={day.getMonth() !== month.getMonth()} aria-label={new Intl.DateTimeFormat(undefined, { dateStyle: 'full' }).format(day)} aria-pressed={date === value} aria-current={(date === today && 'date') || undefined} disabled={unavailable(date)} onClick={() => onChange(date)}>{day.getDate()}</Button>;
        })}
      </Box>
      <Box className={styles.heading}>
        {clearable && <Button variant={BUTTON_VARIANT.OUTLINE} disabled={disabled} onClick={() => onChange(null)}>{__('Clear')}</Button>}
        <Button variant={BUTTON_VARIANT.OUTLINE} disabled={unavailable(today)} onClick={() => onChange(today)}>{__('Today')}</Button>
      </Box>
    </Box>
  );
};
