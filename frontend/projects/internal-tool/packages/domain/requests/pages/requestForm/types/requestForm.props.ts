import type { Field } from '../../../../../business/api/types';

export interface RequestTarget {
  readonly templateId: string;
  readonly environmentId: string;
  readonly schoolId: string;
  readonly schoolVersion: string;
  readonly schoolStatus: 'ACTIVE' | 'DISABLED';
  readonly schoolName: string;
}
export interface RequestFormProps {
  readonly close: () => void;
  readonly success: () => void;
  readonly target?: RequestTarget;
}
export interface SchemaFieldProps {
  readonly name: string;
  readonly field: Field;
}
