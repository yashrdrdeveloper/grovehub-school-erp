import type { ReactElement } from 'react';

import type { InitI18nStatus } from '../../../hooks/useInitI18n/constants/useInitI18n.status';

export interface I18nBootstrapFallbackProps {
  readonly status: InitI18nStatus;
}

export type I18nBootstrapFallbackComponent = (
  props: I18nBootstrapFallbackProps,
) => ReactElement;
