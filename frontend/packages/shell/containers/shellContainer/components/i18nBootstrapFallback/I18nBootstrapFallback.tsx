import { Box, BOX_TAG } from '@Components/components/Box';

import type { I18nBootstrapFallbackComponent } from './types/i18nBootstrapFallback';

/** Must not call __(): the global resolver is not installed yet. */
export const I18nBootstrapFallback: I18nBootstrapFallbackComponent = () => (
  <Box as={BOX_TAG.DIV} role="status" aria-live="polite" dataTestId="i18n-boot-status">
    Loading school workspace…
  </Box>
);
