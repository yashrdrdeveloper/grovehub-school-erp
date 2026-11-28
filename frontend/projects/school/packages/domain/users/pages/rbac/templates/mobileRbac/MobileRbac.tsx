import { RbacContent } from '../../components/rbacContent';
import type { RbacControls } from '../../hooks/useRbacPage/types/useRbacPage';

export const MobileRbac = (props: { controls: RbacControls }) => <RbacContent {...props} />;
