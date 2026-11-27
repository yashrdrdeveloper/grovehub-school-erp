import { Box, BOX_TAG } from '@Components/components/Box';
import { LockKeyhole } from '@Components/icons';

import { LOGIN_LAYOUT_LEGACY_CLASSES } from '../../constants/loginLayout.legacyClasses';
import type { LoginLayoutFooterComponent } from './types/loginLayoutFooter';

/** Auth trust and location footer, stable across sign-in and recovery. */
export const LoginLayoutFooter: LoginLayoutFooterComponent = () => (
  <>
    <Box
      as={BOX_TAG.DIV}
      className={LOGIN_LAYOUT_LEGACY_CLASSES.CLOSING}
      aria-label={__('Secure RDR School workspace')}
    >
      <LockKeyhole size={13} strokeWidth={1.8} aria-hidden="true" />
      <Box as={BOX_TAG.SPAN}>{__('A secure workspace for RDR School')}</Box>
    </Box>
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.FOOTER}>
      <Box as={BOX_TAG.SPAN}>{__('© RDR School')}</Box>
      <Box as={BOX_TAG.SPAN}>{__('Haluhera · Rewari')}</Box>
    </Box>
  </>
);
