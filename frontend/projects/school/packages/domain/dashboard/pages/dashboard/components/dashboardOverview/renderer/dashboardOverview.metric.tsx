import { cn } from '@Base/utilities/cn';
import { Box, BOX_TAG } from '@Components/components/Box';

import { DASHBOARD_TONE_CLASSES } from '../../../constants/dashboard.toneClasses';
import styles from '../DashboardOverview.module.css';
import { getDashboardOverviewUnavailableLabel } from '../helpers/dashboardOverview.getUnavailableLabel';
import type { RenderDashboardOverviewMetric } from './types/dashboardOverview.metric';

/** Pure metric-card renderer, without state or browser effects. */
export const renderDashboardOverviewMetric: RenderDashboardOverviewMetric = ({
  name,
  description,
  icon: Icon,
  tone,
}) => (
  <Box
    as={BOX_TAG.DIV}
    className={cn(styles['school-home__stat'], DASHBOARD_TONE_CLASSES[tone])}
    key={name}
  >
    <Box as={BOX_TAG.DIV} className={styles['school-home__stat-label']}>
      <Box as={BOX_TAG.SPAN}>{name}</Box>
      <Icon size={18} />
    </Box>
    <Box as={BOX_TAG.STRONG} aria-label={getDashboardOverviewUnavailableLabel(name)}>
      —
    </Box>
    <Box as={BOX_TAG.SMALL}>{description}</Box>
  </Box>
);
