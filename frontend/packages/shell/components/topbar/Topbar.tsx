import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { ChevronDown, Menu } from '@Components/icons';
import type { TopbarComponent } from './types/topbar';

export const Topbar: TopbarComponent = ({ accountLabel, schoolLabel, model }) => (
  <Box as="header" className="topbar">
    <Box as="div" className="topbar__left">
      <Button unstyled type="button" className="icon-button topbar__menu" onClick={model.openSidebar}
        aria-label="Open menu"><Menu size={20}/></Button>
      <Box as="span" className="topbar__breadcrumb">{schoolLabel}</Box>
      <Box as="span" className="topbar__separator">/</Box><Box as="span" className="topbar__muted">Workspace</Box>
    </Box>
    <Box as="div" className="topbar__right">
      <Button unstyled type="button" className="icon-button topbar__quick" onClick={model.openAccountPreferences}
        aria-label="Open account preferences">
        <Box as="span" className="profile-avatar profile-avatar--small">{accountLabel.charAt(0).toUpperCase()}</Box>
        <ChevronDown size={14}/>
      </Button>
    </Box>
  </Box>
);
