import { cn } from '@Base/utilities/cn';
import { Box, BOX_TAG } from '@Components/components/Box';
import { ArrowUpRight } from '@Components/icons';

import { DASHBOARD_TONE_CLASSES } from '../../../constants/dashboard.toneClasses';
import styles from '../DashboardMain.module.css';
import type { RenderDashboardMainWorkspace } from './types/dashboardMain.workspace';

/** Pure renderer for the repeated operational previews. */
export const renderDashboardMainWorkspace: RenderDashboardMainWorkspace = ({
  title,
  hint,
  icon: Icon,
  tone,
}) => (
  <Box
    as={BOX_TAG.DIV}
    className={cn(styles['school-home__operation'], DASHBOARD_TONE_CLASSES[tone])}
    key={title}
  >
    <Box as={BOX_TAG.DIV} className={styles['school-home__operation-icon']}>
      <Icon size={18} />
    </Box>
    <Box as={BOX_TAG.DIV}>
      <Box as={BOX_TAG.H3}>{title}</Box>
      <Box as={BOX_TAG.P}>{hint}</Box>
      <Box as={BOX_TAG.SMALL}>{__('Ready for future service data')}</Box>
    </Box>
    <ArrowUpRight size={16} className={styles['school-home__operation-arrow']} />
  </Box>
);
