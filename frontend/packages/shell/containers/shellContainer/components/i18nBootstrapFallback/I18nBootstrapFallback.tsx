import { Box, BOX_TAG } from '@Components/components/Box';

import { InitI18nStatus } from '../../hooks/useInitI18n';
import { I18N_BOOTSTRAP_FALLBACK_TEST_IDS } from './constants/i18nBootstrapFallback.testIds';
import type { I18nBootstrapFallbackComponent } from './types/i18nBootstrapFallback';

/** Must not call __(): the global resolver is not installed yet. */
export const I18nBootstrapFallback: I18nBootstrapFallbackComponent = ({ status }) => {
  const message =
    (status === InitI18nStatus.FAILED && 'Unable to load RDR School. Please refresh the page.') ||
    'Loading school workspace…';

  return (
    <Box
      as={BOX_TAG.DIV}
      role="status"
      aria-live="polite"
      dataTestId={I18N_BOOTSTRAP_FALLBACK_TEST_IDS.STATUS}
    >
      {message}
    </Box>
  );
};
