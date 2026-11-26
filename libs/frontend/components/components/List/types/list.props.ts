import type { ReactNode } from 'react';

export interface ListProps<T> {
  readonly items: readonly T[];
  readonly getItemId: (item: T) => string;
  readonly renderItem: (item: T, index: number) => ReactNode;
  readonly label: string;
  readonly emptyContent?: ReactNode;
  readonly dataTestId?: string;
}
