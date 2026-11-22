import type { ReactNode } from 'react';

import type { RecoveryProps } from '../../types/recovery.props';

export interface RecoveryContentRenderProps {
  readonly props: RecoveryProps;
  readonly formVisible: boolean;
}

export type RenderRecoveryContent = (params: RecoveryContentRenderProps) => ReactNode;
