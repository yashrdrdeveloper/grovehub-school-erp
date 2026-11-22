import { Box } from '@Components/components/Box';
import styles from './DashboardHeader.module.css';
import { CalendarDays, ChevronRight, School } from '@Components/icons';

import type { DashboardHeaderComponent } from './types/dashboardHeader';

/** Presentational greeting and local date. */
export const DashboardHeader: DashboardHeaderComponent = ({ date }) => (
  <Box as="header" className={styles['school-home__header']}>
    <Box as="div">
      <Box as="div" className={styles['school-home__breadcrumb']}>
        <School size={14} /> RDR SCHOOL <ChevronRight size={13} /> OVERVIEW
      </Box>
      <Box as="h1">
        Good {date.greeting}
        <Box as="span" className={styles['school-home__spark']}>
          ✳
        </Box>
      </Box>
      <Box as="p">Your school day, in one clear view.</Box>
    </Box>
    <Box as="div" className={styles['school-home__date']}>
      <CalendarDays size={17} />
      <Box as="span">
        <Box as="strong">{date.dayName}</Box>
        <Box as="small">{date.dateLabel}</Box>
      </Box>
    </Box>
  </Box>
);
