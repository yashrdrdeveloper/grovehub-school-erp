import type { Field } from '../../../../../business/api/types';

export interface RequestFormProps {
  readonly close: () => void;
  readonly success: () => void;
}
export interface SchemaFieldProps {
  readonly name: string;
  readonly field: Field;
}
