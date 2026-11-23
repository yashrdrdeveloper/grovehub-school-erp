import { AcademicPathwaysSvg } from '@Assets/auth';
import { Box, BOX_TAG } from '@Components/components/Box';
import { SVGRenderer } from '@Components/components/SVGRenderer';

import styles from './AuthAmbientGraphic.module.css';
import { AUTH_AMBIENT_TEST_IDS } from './constants/authAmbientGraphic.testIds';

/** Decorative school pathways use a true owner-scoped responsive CSS Module. */
export const AuthAmbientGraphic = () => (
  <Box
    as={BOX_TAG.DIV}
    className={styles.root}
    dataTestId={AUTH_AMBIENT_TEST_IDS.ART}
    aria-hidden="true"
  >
    <SVGRenderer
      component={AcademicPathwaysSvg}
      decorative
      className={styles.svg}
      dataTestId={AUTH_AMBIENT_TEST_IDS.SVG}
    />
    <Box as={BOX_TAG.DIV} className={styles.caption}>
      <Box as={BOX_TAG.SPAN}>{__('LEARNING')}</Box>
      <Box as={BOX_TAG.SPAN} className={styles.separator} />
      <Box as={BOX_TAG.SPAN}>{__('PEOPLE')}</Box>
      <Box as={BOX_TAG.SPAN} className={styles.separator} />
      <Box as={BOX_TAG.SPAN}>{__('PROGRESS')}</Box>
    </Box>
  </Box>
);
