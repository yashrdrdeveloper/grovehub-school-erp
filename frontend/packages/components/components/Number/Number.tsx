import { Input } from '../Input';

import type { NumberProps } from './types/number.props';

export const Number = ({ value, onChange, ...props }: NumberProps) => (
  <Input
    {...props}
    type="number"
    value={value ?? ''}
    onChange={(event) => {
      const next = event.currentTarget.valueAsNumber;
      if (globalThis.Number.isFinite(next)) {
        onChange(next);
        return;
      }
      onChange(null);
    }}
  />
);
