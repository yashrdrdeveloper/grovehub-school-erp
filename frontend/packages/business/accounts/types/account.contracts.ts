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
  pageInfo?: { rows: number; cursor?: string };
}

export interface SearchResult<T> {
  items: T[];
  pageInfo: { rows: number; hasNext: boolean; nextCursor?: string; totalRecords?: number };
}

export type AccountSearchResult = SearchResult<Account>;
