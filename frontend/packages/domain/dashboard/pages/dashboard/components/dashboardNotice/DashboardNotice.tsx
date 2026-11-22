import { Box } from '@Components/components/Box';
import { CircleDashed, Sparkles } from '@Components/icons';

import styles from './DashboardNotice.module.css';

import type { DashboardNoticeComponent } from './types/dashboardNotice';

/** Explicit dashboard-preview disclosure with no invented school activity. */
export const DashboardNotice: DashboardNoticeComponent = () => (
  <Box as="div" className={styles.root} role="status">
    <Box as="div" className={styles.mark}>
      <Sparkles size={17} aria-hidden="true" />
    </Box>

    <Box as="div" className={styles.copy}>
      <Box as="strong">Your workspace is taking shape</Box>
      <Box as="span">
        This is a preview of the dashboard layout. Your school's activity will appear as services
        are connected.
      </Box>
    </Box>

    <Box as="span" className={styles.state}>
      <CircleDashed size={14} aria-hidden="true" /> PREVIEW
    </Box>
  </Box>
);
