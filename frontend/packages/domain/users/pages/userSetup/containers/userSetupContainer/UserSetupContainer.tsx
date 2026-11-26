import { UserSetup } from '../../UserSetup';
import { useUserSetup } from '../../hooks/useUserSetup';

export const UserSetupContainer = () => {
  const controls = useUserSetup();
  return <UserSetup controls={controls} />;
};
