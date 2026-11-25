import { useState } from 'react';
import { CalendarDays } from '@Components/icons';

import { Box } from '../Box';
import { Input } from '../Input';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { BUTTON_VARIANT } from '../Button';
import { Calendar } from '../Calendar';
import { Popover } from '../Popover';

import styles from './Date.module.css';
import type { DateProps } from './types/date.props';

export const Date = ({ value, onChange, dataTestId, min, max, disabled, readOnly, ...props }: DateProps) => {
  const [open, setOpen] = useState(false);
  const minimum = (typeof min === 'string' && min) || undefined;
  const maximum = (typeof max === 'string' && max) || undefined;
  return (
    <Box className={styles.field} dataTestId={dataTestId}>
      <Input {...props} disabled={disabled} readOnly={readOnly} min={min} max={max} type="date" value={value ?? ''} onChange={(event) => onChange(event.currentTarget.value || null)} dataTestId={dataTestId && `${dataTestId}-input`} />
      <Popover title={props.label} open={open} onOpenChange={setOpen} width={330} trigger={<ButtonWithIcon icon={CalendarDays} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Open calendar')} disabled={disabled || readOnly} size={props.controlSize} onClick={() => setOpen(true)} />}>
        <Calendar value={value} min={minimum} max={maximum} clearable={!props.required} onChange={(next) => {
          onChange(next);
          setOpen(false);
        }} />
      </Popover>
    </Box>
  );
};
