import type { Account } from '@Business/domain/users';
import type { FilterCondition } from '@Components/components/Filters';
import type { SortCondition } from '@Components/components/Sort';
import type { UserAction } from '../constants/userSetup.actions';

export interface UserSetupQueryState {
  readonly searchText: string;
  readonly filters: readonly FilterCondition[];
  readonly sorting: readonly SortCondition[];
  readonly quickFilterIds: readonly string[];
  readonly pageSize: number;
  readonly cursors: readonly (string | undefined)[];
}

export interface UserSetupDialog {
  readonly action: UserAction;
  readonly row?: Account;
}
