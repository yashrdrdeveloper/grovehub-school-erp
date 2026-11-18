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
  searchTerms?: string[];
  filters?: Array<{ field: string; operator: string; values: string[] }>;
  sort?: { field: string; direction: 'ASC' | 'DESC' };
  pageInfo?: { rows: number; cursor?: string };
}

export interface SearchResult<T> {
  items: T[];
  pageInfo: { rows: number; hasNext: boolean; nextCursor?: string };
}

export type AccountSearchResult = SearchResult<Account>;
