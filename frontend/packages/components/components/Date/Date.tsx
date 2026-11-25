import { Input } from '../Input';

import type { DateProps } from './types/date.props';

export const Date = ({ value, onChange, dataTestId, ...props }: DateProps) => (
  <Input {...props} dataTestId={dataTestId} type="date" value={value ?? ''} onChange={(event) => onChange(event.currentTarget.value || null)} />
);
