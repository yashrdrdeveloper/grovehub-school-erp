import { Box } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { ChevronDown, Menu } from '@Components/icons';

import { TOPBAR_COPY } from './constants/topbar.copy';
import { TOPBAR_TEST_IDS } from './constants/topbar.testIds';

import type { TopbarComponent } from './types/topbar';

/** Pure workspace chrome, with native semantics delegated to shared components. */
export const Topbar: TopbarComponent = ({ accountLabel, schoolLabel, model }) => (
  <Box as="header" className="topbar" dataTestId={TOPBAR_TEST_IDS.ROOT}>
    <Box as="div" className="topbar__left">
      <Button
        unstyled
        type="button"
        className="icon-button topbar__menu"
        dataTestId={TOPBAR_TEST_IDS.OPEN_MENU}
        onClick={model.openSidebar}
        aria-label={TOPBAR_COPY.OPEN_MENU}
      >
        <Menu size={20} />
      </Button>

      <Box as="span" className="topbar__breadcrumb">
        {schoolLabel}
      </Box>
      <Box as="span" className="topbar__separator" aria-hidden="true">
        /
      </Box>
      <Box as="span" className="topbar__muted">
        {TOPBAR_COPY.WORKSPACE}
      </Box>
    </Box>

    <Box as="div" className="topbar__right">
      <Button
        unstyled
        type="button"
        className="icon-button topbar__quick"
        dataTestId={TOPBAR_TEST_IDS.ACCOUNT}
        onClick={model.openAccountPreferences}
        aria-label={TOPBAR_COPY.ACCOUNT_PREFERENCES}
      >
        <Box as="span" className="profile-avatar profile-avatar--small">
          {accountLabel.charAt(0).toUpperCase()}
        </Box>
        <ChevronDown size={14} aria-hidden="true" />
      </Button>
    </Box>
  </Box>
);
