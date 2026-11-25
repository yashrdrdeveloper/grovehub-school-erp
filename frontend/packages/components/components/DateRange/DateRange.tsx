import { _map } from '@Base/utilities/map';
import { CalendarDays } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../ButtonWithIcon';
import { Calendar } from '../Calendar';
import { Popover } from '../Popover';
import { Tooltip } from '../Tooltip';

import { DATE_RANGE_PRESETS } from './constants/dateRange.presets';
import { shiftRangeMonth } from './helpers/dateRange.calendar';
import { getFormattedRangeDate } from './helpers/dateRange.getFormattedDate';
import { useDateRange } from './hooks/useDateRange';
import styles from './DateRange.module.css';
import type { DateRangeProps } from './types/dateRange.props';

export const DateRange = (props: DateRangeProps) => {
  const { value, disabled, dataTestId } = props;
  const controls = useDateRange(props);
  const empty = !value.startDate && !value.endDate;
  const getSummary = () => {
    if (empty) return __('Select date range');
    return [getFormattedRangeDate(value.startDate), getFormattedRangeDate(value.endDate)].join(' – ');
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
        title={__('Date range')}
        width={724}
        compact
        fullWidth
        trigger={
          <ButtonWithIcon icon={CalendarDays} iconSize={14} iconPosition={BUTTON_ICON_POSITION.RIGHT} variant={BUTTON_VARIANT.OUTLINE} size={BUTTON_SIZE.SMALL} fullWidth className={styles.trigger} data-placeholder={empty} disabled={disabled} aria-label={__('Select date range')} aria-expanded={controls.open} onClick={() => controls.onOpenChange(true)}>
            {getSummary()}
          </ButtonWithIcon>
        }
      >
        <Box className={styles.range}>
          <Box className={styles.layout}>
            <Box className={styles.presets} role="group" aria-label={__('Date range presets')}>
              {_map(DATE_RANGE_PRESETS, (preset) => (
                <Button key={preset.id} variant={BUTTON_VARIANT.GHOST} size={BUTTON_SIZE.SMALL} className={styles.preset} aria-pressed={controls.selectedPreset === preset.id} disabled={disabled} onClick={() => controls.selectPreset(preset)}>
                  {preset.label}
                </Button>
              ))}
            </Box>
            <Box className={styles.selection}>
              <Box className={styles.summary} aria-live="polite">
                <Box className={styles.date} data-placeholder={!controls.draft.startDate}>{getFormattedRangeDate(controls.draft.startDate, __('Start date'))}</Box>
                <Box className={styles.separator} aria-hidden="true">–</Box>
                <Box className={styles.date} data-placeholder={!controls.draft.endDate}>{getFormattedRangeDate(controls.draft.endDate, __('End date'))}</Box>
              </Box>
              <Box className={styles.calendars}>
                <Calendar compact hideFooter navigation="previous" month={controls.month} onMonthChange={controls.setMonth} value={controls.draft.startDate} range={controls.preview} onChange={controls.selectDate} onDayHover={controls.setHover} disabled={disabled} />
                <Calendar compact hideFooter navigation="next" month={controls.nextMonth} onMonthChange={(month) => controls.setMonth(shiftRangeMonth(month, -1))} value={controls.draft.endDate} range={controls.preview} onChange={controls.selectDate} onDayHover={controls.setHover} disabled={disabled} />
              </Box>
            </Box>
          </Box>
          {controls.invalid && <Box role="alert" className={styles.hint}>{__('The start date must not exceed the end date.')}</Box>}
          <Box className={styles.actions}>
            <Box className={styles.hint}>{__('Select a start and end date, or choose a preset.')}</Box>
            <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} disabled={disabled} onClick={controls.clear}>{__('Clear')}</Button>
            <Tooltip label={getApplyReason()} focusable={!controls.canApply}>
              <Button size={BUTTON_SIZE.SMALL} disabled={!controls.canApply} onClick={controls.apply}>{__('Apply')}</Button>
            </Tooltip>
          </Box>
        </Box>
      </Popover>
    </Box>
  );
};
