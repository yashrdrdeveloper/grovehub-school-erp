import { _map } from '@Base/utilities/map';
import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';

import { formatCalendarDate, getCalendarDays } from '../../helpers/calendar.date';
import { getCalendarYears } from '../../helpers/calendar.navigation';
import styles from '../../Calendar.module.css';
import type { CalendarGridProps } from './types/calendarGrid.props';

export const CalendarGrid = ({
  mode,
  month,
  value,
  range,
  disabled,
  size,
  unavailable,
  chooseMonth,
  chooseYear,
  onChange,
  onDayHover,
}: CalendarGridProps) => {
  const year = month.getFullYear();
  const years = getCalendarYears(year);
  const title = new Intl.DateTimeFormat(undefined, { month: 'long', year: 'numeric' }).format(
    month,
  );
  const today = formatCalendarDate(new Date());
  const days = getCalendarDays(month);
  const getPeriodVariant = (selected: boolean) => {
    if (selected) return BUTTON_VARIANT.PRIMARY;
    return BUTTON_VARIANT.GHOST;
  };
  const isSelectedDate = (date: string) => {
    if (range) return date === range.startDate || date === range.endDate;
    return date === value;
  };
  const dayNames = Array.from({ length: 7 }, (_, index) =>
    new Intl.DateTimeFormat(undefined, { weekday: 'short' }).format(new Date(2024, 0, index + 1)),
  );
  return (
    <>
      {mode === 'month' && (
        <Box className={styles.periodGrid} role="group" aria-label={__('Select month')}>
          {_map(
            Array.from({ length: 12 }, (_, index) => new Date(year, index, 1)),
            (entry, index) => (
              <Button
                key={index}
                size={BUTTON_SIZE.SMALL}
                variant={getPeriodVariant(index === month.getMonth())}
                className={styles.periodCell}
                aria-label={new Intl.DateTimeFormat(undefined, {
                  month: 'long',
                  year: 'numeric',
                }).format(entry)}
                aria-pressed={index === month.getMonth()}
                disabled={disabled}
                onClick={() => chooseMonth(index)}
              >
                {new Intl.DateTimeFormat(undefined, { month: 'short' }).format(entry)}
              </Button>
            ),
          )}
        </Box>
      )}
      {mode === 'year' && (
        <Box className={styles.periodGrid} role="group" aria-label={__('Select year')}>
          {_map(years, (entry) => (
            <Button
              key={entry}
              size={BUTTON_SIZE.SMALL}
              variant={getPeriodVariant(entry === year)}
              className={styles.periodCell}
              aria-pressed={entry === year}
              disabled={disabled}
              onClick={() => chooseYear(entry)}
            >
              {entry}
            </Button>
          ))}
        </Box>
      )}
      {mode === 'day' && (
        <Box
          className={styles.days}
          role="group"
          aria-label={title}
          onMouseLeave={() => onDayHover?.(null)}
        >
          {_map(dayNames, (name, index) => (
            <Box key={index} className={styles.weekday}>
              {name}
            </Box>
          ))}
          {_map(days, (day) => {
            const date = formatCalendarDate(day);
            const outside = day.getMonth() !== month.getMonth();
            const selected = isSelectedDate(date);
            const inRange = Boolean(
              range?.startDate && range.endDate && date > range.startDate && date < range.endDate,
            );
            const getVariant = () => {
              if (selected && (!range || !outside)) return BUTTON_VARIANT.PRIMARY;
              return BUTTON_VARIANT.GHOST;
            };
            return (
              <Button
                key={date}
                size={size}
                variant={getVariant()}
                className={styles.day}
                data-outside={outside}
                data-range-day={Boolean(range)}
                data-in-range={inRange && !outside}
                data-range-start={range?.startDate === date && !outside}
                data-range-end={range?.endDate === date && !outside}
                aria-label={new Intl.DateTimeFormat(undefined, { dateStyle: 'full' }).format(day)}
                aria-pressed={selected && (!range || !outside)}
                aria-current={(date === today && 'date') || undefined}
                disabled={unavailable(date) || Boolean(range && outside)}
                onMouseEnter={() => {
                  if (!outside && !unavailable(date)) onDayHover?.(date);
                }}
                onClick={() => onChange(date)}
              >
                {day.getDate()}
              </Button>
            );
          })}
        </Box>
      )}
    </>
  );
};
