import { Box } from '@Components/components/Box';
import { AcademicPathwaysSvg } from '@Assets/auth';
import { AUTH_AMBIENT_TEST_IDS } from './constants/authAmbientGraphic.testIds';
import { SVGRenderer } from '@Components/components/SVGRenderer';

/** Feature-specific frame; SVG shape is owned by @Assets. */
export const AuthAmbientGraphic = () => (
  <Box as="div" className="rdr-auth-art" dataTestId={AUTH_AMBIENT_TEST_IDS.ART} aria-hidden="true">
    <SVGRenderer component={AcademicPathwaysSvg} decorative
      className="rdr-auth-art__svg" dataTestId={AUTH_AMBIENT_TEST_IDS.SVG}/>
    <Box as="div" className="rdr-auth-art__caption">
      <Box as="span">LEARNING</Box><Box as="span" className="rdr-auth-art__sep"/> <Box as="span">PEOPLE</Box><Box as="span" className="rdr-auth-art__sep"/> <Box as="span">PROGRESS</Box>
    </Box>
  </Box>
);
