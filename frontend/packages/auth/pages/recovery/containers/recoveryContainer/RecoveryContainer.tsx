import { Recovery } from '../../Recovery';
import { useRecovery } from '../../hooks/useRecovery';
import type { RecoveryContainerComponent } from './types/recoveryContainer';

export const RecoveryContainer: RecoveryContainerComponent = ({ purpose }) => {
  const controller = useRecovery(purpose);
  return <Recovery {...controller}/>;
};
