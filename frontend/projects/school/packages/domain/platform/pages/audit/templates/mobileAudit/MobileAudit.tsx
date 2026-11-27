import { AuditContent } from '../../components/auditContent';
import type { AuditControls } from '../../hooks/useAuditPage/types/useAuditPage';

export const MobileAudit = (props: { controls: AuditControls }) => <AuditContent {...props} />;
