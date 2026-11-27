import { getRbacTemplate } from './renderer/rbac.template';
import type { RbacProps } from './types/rbac.props';

export const RbacPage = ({ controls, isMobile }: RbacProps) => {
  const Template = getRbacTemplate(isMobile);
  return <Template controls={controls} />;
};
