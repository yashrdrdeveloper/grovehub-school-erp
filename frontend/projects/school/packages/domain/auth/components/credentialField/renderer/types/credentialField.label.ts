import type { ReactNode } from 'react';

import type { CredentialFieldProps } from '../../types/credentialField.props';

export type RenderCredentialFieldLabel = (
  props: Pick<CredentialFieldProps, 'id' | 'label' | 'groupLabel' | 'labelAction'>,
) => ReactNode;
