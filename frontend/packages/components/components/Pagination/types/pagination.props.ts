import type { PaginationType } from '../constants/pagination.type';

export type PaginationAction =
  | { readonly type: 'NEXT' | 'PREVIOUS' | 'LOAD_MORE' }
  | { readonly type: 'PAGE_CHANGE'; readonly page: number }
  | { readonly type: 'PAGE_SIZE_CHANGE'; readonly pageSize: number };

interface PaginationBaseProps {
  readonly pageSize: number;
  readonly pageSizeOptions?: readonly number[];
  readonly disabled?: boolean;
  readonly loading?: boolean;
  readonly onChange: (action: PaginationAction) => void;
  readonly dataTestId?: string;
}

export type PaginationProps = PaginationBaseProps & (
  | { readonly type: typeof PaginationType.CURSOR; readonly currentPage: number; readonly hasNext: boolean; readonly hasPrevious: boolean }
  | { readonly type: typeof PaginationType.OFFSET; readonly currentPage: number; readonly totalRecords: number }
  | { readonly type: typeof PaginationType.INFINITE; readonly hasNext: boolean; readonly loadedRecords: number }
);
