import { getAuditTemplate } from './renderer/audit.template';
import type { AuditProps } from './types/audit.props';

export const AuditPage = ({ controls, isMobile }: AuditProps) => {
  const Template = getAuditTemplate(isMobile);
  return <Template controls={controls} />;
};
