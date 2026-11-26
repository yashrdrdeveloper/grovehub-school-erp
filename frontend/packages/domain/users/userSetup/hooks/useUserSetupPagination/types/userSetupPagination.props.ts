import type { Dispatch, SetStateAction } from 'react';
import type { UserSetupQueryState } from '../../../types/userSetup.state';
import type { useUserSetupQueries } from '../../useUserSetupQueries';

export interface UserSetupPaginationProps {
  readonly query: ReturnType<typeof useUserSetupQueries>['query'];
  readonly setState: Dispatch<SetStateAction<UserSetupQueryState>>;
  readonly resetQuery: (change: Partial<UserSetupQueryState>) => void;
}
