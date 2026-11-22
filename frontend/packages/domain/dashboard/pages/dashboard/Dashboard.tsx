import { Box } from '@Components/components/Box';
import styles from './Dashboard.module.css';

import { DashboardHeader } from './components/dashboardHeader';
import { DASHBOARD_TEST_IDS } from './constants/dashboard.testIds';
import { DashboardNotice } from './components/dashboardNotice';
import { DashboardOverview } from './components/dashboardOverview';
import { DashboardMain } from './components/dashboardMain';
import { DashboardAside } from './components/dashboardAside';
import type { DashboardComponent } from './types';

/** Pure dashboard composition. All data is explicitly marked as a preview. */
export const Dashboard: DashboardComponent = (props) => (
  <Box as="div" className={styles.root} data-testid={DASHBOARD_TEST_IDS.PAGE}>
    <DashboardHeader date={props.date} />
    <DashboardNotice />
    <DashboardOverview />
    <Box as="div" className={styles.layout}>
      <DashboardMain />
      <DashboardAside {...props} />
    </Box>
  </Box>
);
