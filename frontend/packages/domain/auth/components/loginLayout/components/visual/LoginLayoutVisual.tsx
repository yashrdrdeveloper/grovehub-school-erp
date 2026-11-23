import { Box, BOX_TAG } from '@Components/components/Box';
import { SchoolMonogram } from '@Components';

import { AuthAmbientGraphic } from '../../../authAmbientGraphic';
import { LOGIN_LAYOUT_LEGACY_CLASSES } from '../../constants/loginLayout.legacyClasses';
import type { LoginLayoutVisualComponent } from './types/loginLayoutVisual';

/** Shared school identity artwork. No session state, effects or API calls. */
export const LoginLayoutVisual: LoginLayoutVisualComponent = () => (
  <Box
    as={BOX_TAG.SECTION}
    className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL}
    aria-label="RDR School branded gradient background"
  >
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.MOBILE_HEADER}>
      <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.MOBILE_IDENTITY} aria-label="RDR School">
        <SchoolMonogram />
        <Box as={BOX_TAG.SPAN} className={LOGIN_LAYOUT_LEGACY_CLASSES.MOBILE_SCHOOL_NAME}>
          RDR School<Box as={BOX_TAG.SMALL}>Your school workspace</Box>
        </Box>
      </Box>
    </Box>
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_BRAND}>
      <SchoolMonogram />
      <Box as={BOX_TAG.SPAN}>
        <Box as={BOX_TAG.STRONG}>RDR School</Box>
        <Box as={BOX_TAG.SMALL}>Rao Deena Ram Vidya Vihar</Box>
      </Box>
    </Box>
    <AuthAmbientGraphic />
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_COPY}>
      <Box as={BOX_TAG.SPAN} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_KICKER}>
        <Box as={BOX_TAG.SPAN} /> THE RDR SCHOOL WORKSPACE
      </Box>
      <Box as={BOX_TAG.H2}>
        For every school day.
        <br />
        <Box as={BOX_TAG.EM}>And everyone behind it.</Box>
      </Box>
      <Box as={BOX_TAG.P}>A connected space for the people who make learning happen.</Box>
      <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_FOOTER}>
        <Box as={BOX_TAG.SPAN}>LEARN</Box>
        <Box as={BOX_TAG.I} aria-hidden="true" />
        <Box as={BOX_TAG.SPAN}>ORGANISE</Box>
        <Box as={BOX_TAG.I} aria-hidden="true" />
        <Box as={BOX_TAG.SPAN}>GROW</Box>
      </Box>
    </Box>
  </Box>
);
