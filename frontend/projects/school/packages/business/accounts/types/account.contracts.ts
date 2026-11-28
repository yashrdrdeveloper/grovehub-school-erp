import type { Account } from '@Business/auth';

export interface PersonNames {
  firstName: string;
  middleName?: string | null;
  lastName: string;
}

export interface CreateAccount {
  person: PersonNames;
  email: string;
  phoneNumber?: string;
}

export interface SearchRequest {
  searchText?: string;
  filters?: Array<{ field: string; operator: string; values: string[] }>;
  sort?: Array<{ field: string; direction: 'ASC' | 'DESC' }>;
  pageInfo?: { pageNumber?: number; rowSize?: number; rows?: number; cursor?: string };
}

export interface SearchResult<T> {
  items: T[];
  pageInfo: {
    rows: number;
    pageNumber?: number;
    rowSize?: number;
    totalPages?: number;
    hasPrevious?: boolean;
    hasNext: boolean;
    nextCursor?: string;
    totalRecords?: number;
  };
}

export type AccountSearchResult = SearchResult<Account>;
