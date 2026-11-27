import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './DashboardOverview.module.css';
import { renderDashboardOverviewMetric } from './renderer/dashboardOverview.metric';
import sectionStyles from '../../DashboardSection.module.css';
import { DASHBOARD_METRICS } from '../../constants/dashboard.metrics';
import type { DashboardOverviewComponent } from './types/dashboardOverview';

export const DashboardOverview: DashboardOverviewComponent = () => (
  <Box
    as={BOX_TAG.SECTION}
    className={sectionStyles['school-home__section']}
    aria-labelledby="home-priorities"
  >
    <Box as={BOX_TAG.DIV} className={sectionStyles['school-home__section-heading']}>
      <Box as={BOX_TAG.DIV}>
        <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__eyebrow']}>
          {__('YOUR DAY / 01')}
        </Box>
        <Box as={BOX_TAG.H2} id="home-priorities">
          {__('At a glance')}
        </Box>
      </Box>
      <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__hint']}>
        {__('No live data yet')}
      </Box>
    </Box>
    <Box as={BOX_TAG.DIV} className={styles['school-home__stats']}>
      {_map(DASHBOARD_METRICS, renderDashboardOverviewMetric)}
    </Box>
  </Box>
);
