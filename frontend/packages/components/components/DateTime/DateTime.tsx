import { Input } from '../Input';

import type { DateTimeProps } from './types/dateTime.props';

export const DateTime = ({ value, onChange, ...props }: DateTimeProps) => (
  <Input {...props} type="datetime-local" value={value ?? ''} onChange={(event) => onChange(event.currentTarget.value || null)} />
);
