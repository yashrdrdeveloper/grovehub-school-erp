import { Box, BOX_TAG } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Menu } from '@Components/icons';
import { SchoolMonogram } from '@Components';

import styles from './Topbar.module.css';

import { TOPBAR_COPY } from './constants/topbar.copy';
import { TOPBAR_TEST_IDS } from './constants/topbar.testIds';

import type { TopbarComponent } from './types/topbar';

/** Pure workspace chrome, with native semantics delegated to shared components. */
export const Topbar: TopbarComponent = ({ schoolLabel, model }) => (
  <Box as={BOX_TAG.HEADER} className={styles.root} dataTestId={TOPBAR_TEST_IDS.ROOT}>
    <Box as={BOX_TAG.DIV} className={styles.left}>
      <Button
        unstyled
        type="button"
        className={styles.menu}
        dataTestId={TOPBAR_TEST_IDS.OPEN_MENU}
        onClick={model.openSidebar}
        aria-label={TOPBAR_COPY.OPEN_MENU}
      >
        <Menu size={20} />
      </Button>

      <Box as={BOX_TAG.SPAN} role="img" aria-label={schoolLabel}>
        <SchoolMonogram className={styles.logo} />
      </Box>
    </Box>
  </Box>
);
