import { Input } from '../Input';

import { getLocalDateTimeValue, getUTCDateTimeValue } from './helpers/dateTime.value';
import type { DateTimeProps } from './types/dateTime.props';

export const DateTime = ({ value, onChange, dataTestId, min, max, ...props }: DateTimeProps) => (
  <Input
    {...props}
    dataTestId={dataTestId}
    type="datetime-local"
    min={getLocalDateTimeValue(String(min ?? '')) || undefined}
    max={getLocalDateTimeValue(String(max ?? '')) || undefined}
    value={getLocalDateTimeValue(value)}
    onChange={(event) => onChange(getUTCDateTimeValue(event.currentTarget.value))}
  />
);
