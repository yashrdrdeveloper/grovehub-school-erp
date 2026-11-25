import type { useUserSetup } from '../../../hooks/useUserSetup';

export interface UserTableProps {
  readonly controls: ReturnType<typeof useUserSetup>;
}
