import { Box, BOX_TAG } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { ChevronDown, Menu } from '@Components/icons';

import styles from './Topbar.module.css';

import { TOPBAR_COPY } from './constants/topbar.copy';
import { TOPBAR_TEST_IDS } from './constants/topbar.testIds';

import type { TopbarComponent } from './types/topbar';

/** Pure workspace chrome, with native semantics delegated to shared components. */
export const Topbar: TopbarComponent = ({ accountLabel, schoolLabel, model }) => (
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

      <Box as={BOX_TAG.SPAN} className={styles.breadcrumb}>
        {schoolLabel}
      </Box>
      <Box as={BOX_TAG.SPAN} className={styles.separator} aria-hidden="true">
        /
      </Box>
      <Box as={BOX_TAG.SPAN} className={styles.muted}>
        {TOPBAR_COPY.WORKSPACE}
      </Box>
    </Box>

    <Box as={BOX_TAG.DIV} className={styles.right}>
      <Button
        unstyled
        type="button"
        className={styles.quick}
        dataTestId={TOPBAR_TEST_IDS.ACCOUNT}
        onClick={model.openAccountPreferences}
        aria-label={TOPBAR_COPY.ACCOUNT_PREFERENCES}
      >
        <Box as={BOX_TAG.SPAN} className={styles.avatar}>
          {accountLabel.charAt(0).toUpperCase()}
        </Box>
        <ChevronDown size={14} aria-hidden="true" />
      </Button>
    </Box>
  </Box>
);
