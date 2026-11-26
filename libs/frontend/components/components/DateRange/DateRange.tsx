import { CalendarDays } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../ButtonWithIcon';
import { Calendar } from '../Calendar';
import { Popover } from '../Popover';
import { Tooltip } from '../Tooltip';

import { CUSTOM_DATE_RANGE_PRESET } from './constants/dateRange.presets';
import { DateRangePresets } from './components/dateRangePresets';
import { shiftRangeMonth } from './helpers/dateRange.calendar';
import { getFormattedRangeDate } from './helpers/dateRange.getFormattedDate';
import { useDateRange } from './hooks/useDateRange';
import styles from './DateRange.module.css';
import type { DateRangeProps } from './types/dateRange.props';

export const DateRange = (props: DateRangeProps) => {
  const { value, disabled, dataTestId, showHeader = true, showFooter = true } = props;
  const controls = useDateRange(props);
  const empty = !value.startDate && !value.endDate;
  const getSummary = () => {
    if (empty) return __('Select date range');
    return [getFormattedRangeDate(value.startDate), getFormattedRangeDate(value.endDate)].join(
      ' – ',
    );
  };
  const getApplyReason = () => {
    if (disabled) return __('Date range is disabled.');
    if (controls.invalid) return __('The start date must not exceed the end date.');
    if (controls.incomplete) return __('Select both a start and an end date.');
    return undefined;
  };

  return (
    <Box dataTestId={dataTestId} className={styles.root}>
      <Popover
        open={controls.open}
        onOpenChange={controls.onOpenChange}
        onOutsideClick={controls.onOutsideClick}
        title={__('Date range')}
        showHeader={showHeader}
        headerContent={
          <Box className={styles.headerSummary} aria-live="polite">
            {!controls.draft.startDate && __('Select date range')}
            {controls.draft.startDate && getFormattedRangeDate(controls.draft.startDate)}
            {controls.draft.startDate && ' – '}
            {controls.draft.startDate &&
              getFormattedRangeDate(controls.draft.endDate, __('End date'))}
          </Box>
        }
        width={520}
        compact
        fullWidth
        trigger={
          <ButtonWithIcon
            icon={CalendarDays}
            iconSize={14}
            iconPosition={BUTTON_ICON_POSITION.RIGHT}
            variant={BUTTON_VARIANT.OUTLINE}
            size={BUTTON_SIZE.SMALL}
            fullWidth
            className={styles.trigger}
            data-placeholder={empty}
            disabled={disabled}
            aria-label={__('Select date range')}
            aria-expanded={controls.open}
            onClick={() => controls.onOpenChange(true)}
          >
            {getSummary()}
          </ButtonWithIcon>
        }
      >
        <Box className={styles.range}>
          <Box className={styles.layout}>
            <DateRangePresets controls={controls} disabled={disabled} />
            <Box className={styles.selection}>
              <Box
                className={styles.calendars}
                key={controls.selectedPreset ?? CUSTOM_DATE_RANGE_PRESET}
              >
                <Calendar
                  compact
                  hideFooter
                  navigation="previous"
                  month={controls.month}
                  onMonthChange={controls.setMonth}
                  value={controls.draft.startDate}
                  range={controls.preview}
                  onChange={controls.selectDate}
                  onDayHover={controls.setHover}
                  disabled={disabled}
                />
                <Calendar
                  compact
                  hideFooter
                  navigation="next"
                  month={controls.nextMonth}
                  onMonthChange={(month) => controls.setMonth(shiftRangeMonth(month, -1))}
                  value={controls.draft.endDate}
                  range={controls.preview}
                  onChange={controls.selectDate}
                  onDayHover={controls.setHover}
                  disabled={disabled}
                />
              </Box>
            </Box>
          </Box>
          {controls.invalid && (
            <Box role="alert" className={styles.hint}>
              {__('The start date must not exceed the end date.')}
            </Box>
          )}
          {showFooter && (
            <Box className={styles.footer}>
              <Button
                size={BUTTON_SIZE.SMALL}
                variant={BUTTON_VARIANT.GHOST}
                disabled={disabled}
                onClick={controls.clear}
              >
                {__('Clear')}
              </Button>
              <Tooltip label={getApplyReason()} focusable={!controls.canApply}>
                <Button
                  size={BUTTON_SIZE.SMALL}
                  disabled={!controls.canApply}
                  onClick={controls.apply}
                >
                  {__('Apply')}
                </Button>
              </Tooltip>
            </Box>
          )}
        </Box>
      </Popover>
    </Box>
  );
};
