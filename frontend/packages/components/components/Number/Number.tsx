import { Input } from '../Input';

import type { NumberProps } from './types/number.props';

export const Number = ({ value, onChange, dataTestId, ...props }: NumberProps) => (
  <Input
    {...props}
    dataTestId={dataTestId}
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
