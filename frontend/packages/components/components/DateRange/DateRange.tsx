import { useState } from 'react';
import { CalendarDays } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../ButtonWithIcon';
import { Popover } from '../Popover';
import { Date } from '../Date';
import { Tooltip } from '../Tooltip';

import styles from './DateRange.module.css';
import type { DateRangeProps } from './types/dateRange.props';

export const DateRange = ({ value, onChange, disabled, dataTestId }: DateRangeProps) => {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState(value);
  const label = __('Date range');
  const empty = value.startDate === null && value.endDate === null;
  const invalid = draft.startDate !== null && draft.endDate !== null && draft.startDate > draft.endDate;
  const summary = () => {
    if (empty) return __('Select date range');
    return [value.startDate ?? '—', value.endDate ?? '—'].join(' – ');
  };
  const onOpenChange = (next: boolean) => {
    if (next) setDraft(value);
    setOpen(next);
  };
  return (
    <Box dataTestId={dataTestId}>
      <Popover open={open} onOpenChange={onOpenChange} title={label} fullWidth trigger={<ButtonWithIcon icon={CalendarDays} iconPosition={BUTTON_ICON_POSITION.RIGHT} variant={BUTTON_VARIANT.OUTLINE} fullWidth className={styles.trigger} data-placeholder={empty} disabled={disabled} onClick={() => onOpenChange(true)}>{summary()}</ButtonWithIcon>}>
        <Box className={styles.range}>
          <Date label={__('Start date')} value={draft.startDate} disabled={disabled} onChange={(next) => setDraft({ ...draft, startDate: next })} />
          <Date label={__('End date')} value={draft.endDate} disabled={disabled} onChange={(next) => setDraft({ ...draft, endDate: next })} />
          {invalid && <Box role="alert">{__('The start must not exceed the end.')}</Box>}
          <Box className={styles.actions}>
            <Button variant={BUTTON_VARIANT.OUTLINE} onClick={() => setDraft({ startDate: null, endDate: null })}>{__('Clear')}</Button>
            <Tooltip label={invalid && __('The start date must not exceed the end date.') || undefined} focusable={invalid}><Button disabled={invalid} onClick={() => {
              onChange(draft);
              setOpen(false);
            }}>{__('Apply')}</Button></Tooltip>
          </Box>
        </Box>
      </Popover>
    </Box>
  );
};
