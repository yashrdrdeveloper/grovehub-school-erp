import { RbacContent } from '../../components/rbacContent';
import type { RbacControls } from '../../hooks/useRbacPage/types/useRbacPage';

export const DesktopRbac = (props: { controls: RbacControls }) => <RbacContent {...props} />;
