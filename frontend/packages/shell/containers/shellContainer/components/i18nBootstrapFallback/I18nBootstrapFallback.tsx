import { Box, BOX_TAG } from '@Components/components/Box';

import { I18N_BOOTSTRAP_FALLBACK_TEST_IDS } from './constants/i18nBootstrapFallback.testIds';
import type { I18nBootstrapFallbackComponent } from './types/i18nBootstrapFallback';

/** Must not call __(): the global resolver is not installed yet. */
export const I18nBootstrapFallback: I18nBootstrapFallbackComponent = () => (
  <Box as={BOX_TAG.DIV} role="status" aria-live="polite" dataTestId={I18N_BOOTSTRAP_FALLBACK_TEST_IDS.STATUS}>
    Loading school workspace…
  </Box>
);
