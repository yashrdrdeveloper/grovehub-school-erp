import { Input } from '../Input';

import type { InputProps } from '../Input';

export const DebouncedInputView = ({ dataTestId, ...props }: InputProps) => <Input {...props} dataTestId={dataTestId} />;
