import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { RbacPage } from '../../RbacPage';
import { useRbacPage } from '../../hooks/useRbacPage';

export const RbacContainer = () => {
  const controls = useRbacPage();
  const isMobile = useMobileViewport();
  return <RbacPage controls={controls} isMobile={isMobile} />;
};
