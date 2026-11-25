import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { ArrowLeft, ArrowRight } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';

import { formatCalendarDate, getCalendarDays, parseCalendarDate } from './helpers/calendar.date';
import { getCalendarYears, shiftCalendarView } from './helpers/calendar.navigation';
import type { CalendarMode } from './types/calendar.mode';
import styles from './Calendar.module.css';
import type { CalendarProps } from './types/calendar.props';

export const Calendar = ({
  value, onChange, min, max, disabled, clearable = true, dataTestId, month: controlledMonth,
  onMonthChange, range, onDayHover, compact, hideFooter, navigation = 'both',
}: CalendarProps) => {
  const [localMonth, setLocalMonth] = useState(() => parseCalendarDate(value) ?? new Date());
  const [mode, setMode] = useState<CalendarMode>('day');
  const month = controlledMonth ?? localMonth;
  const year = month.getFullYear();
  const years = getCalendarYears(year);
  const monthLabel = new Intl.DateTimeFormat(undefined, { month: 'long' }).format(month);
  const title = new Intl.DateTimeFormat(undefined, { month: 'long', year: 'numeric' }).format(month);
  const today = formatCalendarDate(new Date());
  const days = getCalendarDays(month);
  const changeMonth = (next: Date) => {
    onDayHover?.(null);
    if (onMonthChange) onMonthChange(next);
    else setLocalMonth(next);
  };
  const shift = (delta: number) => changeMonth(shiftCalendarView(month, mode, delta));
  const chooseMonth = (index: number) => {
    changeMonth(new Date(year, index, 1));
    setMode('day');
  };
  const chooseYear = (next: number) => {
    changeMonth(new Date(next, month.getMonth(), 1));
    setMode('month');
  };
  const getNavigationLabel = (previous: boolean) => {
    if (mode === 'year') return previous ? __('Previous decade') : __('Next decade');
    if (mode === 'month') return previous ? __('Previous year') : __('Next year');
    return previous ? __('Previous month') : __('Next month');
  };
  const unavailable = (date: string) => Boolean(disabled || (min && date < min) || (max && date > max));
  const dayNames = Array.from({ length: 7 }, (_, index) => new Intl.DateTimeFormat(undefined, { weekday: 'short' }).format(new Date(2024, 0, index + 1)));
  return (
    <Box className={styles.calendar} data-compact={compact} dataTestId={dataTestId}>
      <Box className={styles.heading}>
        <ButtonWithIcon icon={ArrowLeft} size={compact ? BUTTON_SIZE.SMALL : BUTTON_SIZE.MEDIUM} variant={BUTTON_VARIANT.GHOST} className={styles.navigation} data-hidden={mode === 'day' && navigation === 'next'} disabled={disabled} aria-label={getNavigationLabel(true)} onClick={() => shift(-1)} />
        <Box className={styles.period} aria-live="polite">
          {mode === 'day' && <Button variant={BUTTON_VARIANT.GHOST} size={BUTTON_SIZE.SMALL} className={styles.periodButton} disabled={disabled} aria-label={__('Choose month')} onClick={() => setMode('month')}>{monthLabel}</Button>}
          {mode !== 'year' && <Button variant={BUTTON_VARIANT.GHOST} size={BUTTON_SIZE.SMALL} className={styles.periodButton} disabled={disabled} aria-label={__('Choose year')} onClick={() => setMode('year')}>{year}</Button>}
          {mode === 'year' && <Box>{years[1]} – {years[10]}</Box>}
        </Box>
        <ButtonWithIcon icon={ArrowRight} size={compact ? BUTTON_SIZE.SMALL : BUTTON_SIZE.MEDIUM} variant={BUTTON_VARIANT.GHOST} className={styles.navigation} data-hidden={mode === 'day' && navigation === 'previous'} disabled={disabled} aria-label={getNavigationLabel(false)} onClick={() => shift(1)} />
      </Box>
      {mode === 'month' && (
        <Box className={styles.periodGrid} role="group" aria-label={__('Select month')}>
          {_map(Array.from({ length: 12 }, (_, index) => new Date(year, index, 1)), (entry, index) => (
            <Button key={index} size={BUTTON_SIZE.SMALL} variant={index === month.getMonth() ? BUTTON_VARIANT.PRIMARY : BUTTON_VARIANT.GHOST} className={styles.periodCell} aria-label={new Intl.DateTimeFormat(undefined, { month: 'long', year: 'numeric' }).format(entry)} aria-pressed={index === month.getMonth()} disabled={disabled} onClick={() => chooseMonth(index)}>
              {new Intl.DateTimeFormat(undefined, { month: 'short' }).format(entry)}
            </Button>
          ))}
        </Box>
      )}
      {mode === 'year' && (
        <Box className={styles.periodGrid} role="group" aria-label={__('Select year')}>
          {_map(years, (entry) => (
            <Button key={entry} size={BUTTON_SIZE.SMALL} variant={entry === year ? BUTTON_VARIANT.PRIMARY : BUTTON_VARIANT.GHOST} className={styles.periodCell} aria-pressed={entry === year} disabled={disabled} onClick={() => chooseYear(entry)}>{entry}</Button>
          ))}
        </Box>
      )}
      {mode === 'day' && <Box className={styles.days} role="group" aria-label={title} onMouseLeave={() => onDayHover?.(null)}>
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
      </Box>}
      {!hideFooter && (
        <Box className={styles.heading}>
          {clearable && <Button variant={BUTTON_VARIANT.OUTLINE} disabled={disabled} onClick={() => onChange(null)}>{__('Clear')}</Button>}
          <Button variant={BUTTON_VARIANT.OUTLINE} disabled={unavailable(today)} onClick={() => onChange(today)}>{__('Today')}</Button>
        </Box>
      )}
    </Box>
  );
};
