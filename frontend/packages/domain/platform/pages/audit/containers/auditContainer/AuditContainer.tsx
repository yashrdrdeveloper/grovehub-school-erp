import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { AuditPage } from '../../AuditPage';
import { useAuditPage } from '../../hooks/useAuditPage';

export const AuditContainer = () => {
  const controls = useAuditPage();
  const isMobile = useMobileViewport();
  return <AuditPage controls={controls} isMobile={isMobile} />;
};
