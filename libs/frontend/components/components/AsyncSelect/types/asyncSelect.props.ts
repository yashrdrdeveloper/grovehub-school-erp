import type { SelectOption, SelectProps, SelectValue } from '../../Select';

export interface AsyncOptionsRequest {
  readonly searchText: string;
  readonly cursor?: string;
  readonly limit: number;
  readonly signal?: AbortSignal;
}

export interface AsyncOptionsResponse<T extends SelectValue> {
  readonly options: readonly SelectOption<T>[];
  readonly hasNext: boolean;
  readonly nextCursor?: string;
}

export interface AsyncSelectSource<T extends SelectValue> {
  readonly cacheKey: readonly string[];
  readonly loadOptions: (request: AsyncOptionsRequest) => Promise<AsyncOptionsResponse<T>>;
  readonly resolveOptions?: (
    values: readonly T[],
    signal?: AbortSignal,
  ) => Promise<readonly SelectOption<T>[]>;
  readonly initialOptions?: readonly SelectOption<T>[];
  readonly pageSize?: number;
  readonly debounceMs?: number;
}

export interface AsyncSelectProps<T extends SelectValue = string>
  extends Omit<SelectProps<T>, 'options'>,
    AsyncSelectSource<T> {
  readonly dataTestId?: string;
}
