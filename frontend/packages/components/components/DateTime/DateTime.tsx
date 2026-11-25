import { Input } from '../Input';

import type { DateTimeProps } from './types/dateTime.props';

export const DateTime = ({ value, onChange, dataTestId, ...props }: DateTimeProps) => (
  <Input {...props} dataTestId={dataTestId} type="datetime-local" value={value ?? ''} onChange={(event) => onChange(event.currentTarget.value || null)} />
);
