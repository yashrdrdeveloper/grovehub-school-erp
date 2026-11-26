import { useState } from 'react';
import { ArrowLeft, ArrowRight } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';

import { formatCalendarDate, parseCalendarDate } from './helpers/calendar.date';
import { getCalendarYears, shiftCalendarView } from './helpers/calendar.navigation';
import { CalendarGrid } from './components/calendarGrid';
import type { CalendarMode } from './types/calendar.mode';
import styles from './Calendar.module.css';
import type { CalendarProps } from './types/calendar.props';

export const Calendar = (props: CalendarProps) => {
  const {
    value,
    onChange,
    min,
    max,
    disabled,
    clearable = true,
    dataTestId,
    month: controlledMonth,
    onMonthChange,
    onDayHover,
    compact,
    hideFooter,
    navigation = 'both',
  } = props;
  const [localMonth, setLocalMonth] = useState(() => parseCalendarDate(value) ?? new Date());
  const [mode, setMode] = useState<CalendarMode>('day');
  const month = controlledMonth ?? localMonth;
  const year = month.getFullYear();
  const years = getCalendarYears(year);
  const monthLabel = new Intl.DateTimeFormat(undefined, { month: 'long' }).format(month);
  const today = formatCalendarDate(new Date());
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
    if (previous) {
      if (mode === 'year') return __('Previous decade');
      if (mode === 'month') return __('Previous year');
      return __('Previous month');
    }
    if (mode === 'year') return __('Next decade');
    if (mode === 'month') return __('Next year');
    return __('Next month');
  };
  const getControlSize = () => {
    if (compact) return BUTTON_SIZE.SMALL;
    return BUTTON_SIZE.MEDIUM;
  };
  const unavailable = (date: string) =>
    Boolean(disabled || (min && date < min) || (max && date > max));
  return (
    <Box className={styles.calendar} data-compact={compact} dataTestId={dataTestId}>
      <Box className={styles.heading}>
        <ButtonWithIcon
          icon={ArrowLeft}
          size={getControlSize()}
          variant={BUTTON_VARIANT.GHOST}
          className={styles.navigation}
          data-hidden={mode === 'day' && navigation === 'next'}
          disabled={disabled}
          aria-label={getNavigationLabel(true)}
          onClick={() => shift(-1)}
        />
        <Box className={styles.period} aria-live="polite">
          {mode === 'day' && (
            <Button
              variant={BUTTON_VARIANT.GHOST}
              size={BUTTON_SIZE.SMALL}
              className={styles.periodButton}
              disabled={disabled}
              aria-label={__('Choose month')}
              onClick={() => setMode('month')}
            >
              {monthLabel}
            </Button>
          )}
          {mode !== 'year' && (
            <Button
              variant={BUTTON_VARIANT.GHOST}
              size={BUTTON_SIZE.SMALL}
              className={styles.periodButton}
              disabled={disabled}
              aria-label={__('Choose year')}
              onClick={() => setMode('year')}
            >
              {year}
            </Button>
          )}
          {mode === 'year' && (
            <Box>
              {years[1]} – {years[10]}
            </Box>
          )}
        </Box>
        <ButtonWithIcon
          icon={ArrowRight}
          size={getControlSize()}
          variant={BUTTON_VARIANT.GHOST}
          className={styles.navigation}
          data-hidden={mode === 'day' && navigation === 'previous'}
          disabled={disabled}
          aria-label={getNavigationLabel(false)}
          onClick={() => shift(1)}
        />
      </Box>
      <CalendarGrid
        {...props}
        mode={mode}
        month={month}
        size={getControlSize()}
        unavailable={unavailable}
        chooseMonth={chooseMonth}
        chooseYear={chooseYear}
      />
      {!hideFooter && (
        <Box className={styles.heading}>
          {clearable && (
            <Button
              variant={BUTTON_VARIANT.OUTLINE}
              disabled={disabled}
              onClick={() => onChange(null)}
            >
              {__('Clear')}
            </Button>
          )}
          <Button
            variant={BUTTON_VARIANT.OUTLINE}
            disabled={unavailable(today)}
            onClick={() => onChange(today)}
          >
            {__('Today')}
          </Button>
        </Box>
      )}
    </Box>
  );
};
