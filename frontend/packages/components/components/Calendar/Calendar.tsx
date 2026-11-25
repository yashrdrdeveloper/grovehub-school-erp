import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { ArrowLeft, ArrowRight } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';

import { formatCalendarDate, getCalendarDays, parseCalendarDate } from './helpers/calendar.date';
import styles from './Calendar.module.css';
import type { CalendarProps } from './types/calendar.props';

export const Calendar = ({
  value, onChange, min, max, disabled, clearable = true, dataTestId, month: controlledMonth,
  onMonthChange, range, onDayHover, compact, hideFooter, navigation = 'both',
}: CalendarProps) => {
  const [localMonth, setLocalMonth] = useState(() => parseCalendarDate(value) ?? new Date());
  const month = controlledMonth ?? localMonth;
  const title = new Intl.DateTimeFormat(undefined, { month: 'long', year: 'numeric' }).format(month);
  const today = formatCalendarDate(new Date());
  const days = getCalendarDays(month);
  const shift = (delta: number) => {
    const next = new Date(month.getFullYear(), month.getMonth() + delta, 1);
    if (onMonthChange) onMonthChange(next);
    else setLocalMonth(next);
  };
  const unavailable = (date: string) => Boolean(disabled || (min && date < min) || (max && date > max));
  const dayNames = Array.from({ length: 7 }, (_, index) => new Intl.DateTimeFormat(undefined, { weekday: 'short' }).format(new Date(2024, 0, index + 1)));
  return (
    <Box className={styles.calendar} data-compact={compact} dataTestId={dataTestId}>
      <Box className={styles.heading}>
        <ButtonWithIcon icon={ArrowLeft} size={compact ? BUTTON_SIZE.SMALL : BUTTON_SIZE.MEDIUM} variant={BUTTON_VARIANT.GHOST} className={styles.navigation} data-hidden={navigation === 'next'} disabled={disabled} aria-label={__('Previous month')} onClick={() => shift(-1)} />
        <Box aria-live="polite">{title}</Box>
        <ButtonWithIcon icon={ArrowRight} size={compact ? BUTTON_SIZE.SMALL : BUTTON_SIZE.MEDIUM} variant={BUTTON_VARIANT.GHOST} className={styles.navigation} data-hidden={navigation === 'previous'} disabled={disabled} aria-label={__('Next month')} onClick={() => shift(1)} />
      </Box>
      <Box className={styles.days} role="group" aria-label={title} onMouseLeave={() => onDayHover?.(null)}>
        {_map(dayNames, (name, index) => <Box key={index} className={styles.weekday}>{name}</Box>)}
        {_map(days, (day) => {
          const date = formatCalendarDate(day);
          const outside = day.getMonth() !== month.getMonth();
          const selected = range ? date === range.startDate || date === range.endDate : date === value;
          const inRange = Boolean(range?.startDate && range.endDate && date > range.startDate && date < range.endDate);
          const getVariant = () => {
            if (selected && (!range || !outside)) return BUTTON_VARIANT.PRIMARY;
            return BUTTON_VARIANT.GHOST;
          };
          return <Button key={date} size={compact ? BUTTON_SIZE.SMALL : BUTTON_SIZE.MEDIUM} variant={getVariant()} className={styles.day} data-outside={outside} data-range-day={Boolean(range)} data-in-range={inRange && !outside} data-range-start={range?.startDate === date && !outside} data-range-end={range?.endDate === date && !outside} aria-label={new Intl.DateTimeFormat(undefined, { dateStyle: 'full' }).format(day)} aria-pressed={selected && (!range || !outside)} aria-current={(date === today && 'date') || undefined} disabled={unavailable(date) || Boolean(range && outside)} onMouseEnter={() => { if (!outside && !unavailable(date)) onDayHover?.(date); }} onClick={() => onChange(date)}>{day.getDate()}</Button>;
        })}
      </Box>
      {!hideFooter && (
        <Box className={styles.heading}>
          {clearable && <Button variant={BUTTON_VARIANT.OUTLINE} disabled={disabled} onClick={() => onChange(null)}>{__('Clear')}</Button>}
          <Button variant={BUTTON_VARIANT.OUTLINE} disabled={unavailable(today)} onClick={() => onChange(today)}>{__('Today')}</Button>
        </Box>
      )}
    </Box>
  );
};
