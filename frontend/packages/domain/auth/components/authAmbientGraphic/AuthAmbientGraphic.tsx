import { AcademicPathwaysSvg } from '@Assets/auth';
import { Box } from '@Components/components/Box';
import { SVGRenderer } from '@Components/components/SVGRenderer';

import styles from './AuthAmbientGraphic.module.css';
import { AUTH_AMBIENT_TEST_IDS } from './constants/authAmbientGraphic.testIds';

/** Decorative school pathways use a true owner-scoped responsive CSS Module. */
export const AuthAmbientGraphic = () => (
  <Box
    as="div"
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
    <Box as="div" className={styles.caption}>
      <Box as="span">LEARNING</Box>
      <Box as="span" className={styles.separator} />
      <Box as="span">PEOPLE</Box>
      <Box as="span" className={styles.separator} />
      <Box as="span">PROGRESS</Box>
    </Box>
  </Box>
);
