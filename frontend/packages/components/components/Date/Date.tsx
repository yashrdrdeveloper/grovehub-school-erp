import { Input } from '../Input';

import type { DateProps } from './types/date.props';

export const Date = ({ value, onChange, ...props }: DateProps) => (
  <Input {...props} type="date" value={value ?? ''} onChange={(event) => onChange(event.currentTarget.value || null)} />
);
