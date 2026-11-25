import type { useUserSetup } from '../hooks/useUserSetup';

export interface UserSetupProps {
  readonly controls: ReturnType<typeof useUserSetup>;
}
