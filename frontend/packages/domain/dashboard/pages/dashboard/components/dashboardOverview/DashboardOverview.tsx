import { _map } from '@Base/utilities/map';
import { cn } from '@Base/utilities/cn';
import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './DashboardOverview.module.css';
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
      {_map(DASHBOARD_METRICS, ({ name, description, icon: Icon, tone }) => (
        <Box
          as={BOX_TAG.DIV}
          className={cn(styles['school-home__stat'], sectionStyles['school-home__tone--' + tone])}
          key={name}
        >
          <Box as={BOX_TAG.DIV} className={styles['school-home__stat-label']}>
            <Box as={BOX_TAG.SPAN}>{name}</Box>
            <Icon size={18} />
          </Box>
          <Box as={BOX_TAG.STRONG} aria-label={name + ' not available'}>
            —
          </Box>
          <Box as={BOX_TAG.SMALL}>{description}</Box>
        </Box>
      ))}
    </Box>
  </Box>
);
