import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './DashboardHeader.module.css';
import { CalendarDays, ChevronRight, School } from '@Components/icons';

import type { DashboardHeaderComponent } from './types/dashboardHeader';

/** Presentational greeting and local date. */
export const DashboardHeader: DashboardHeaderComponent = ({ date }) => (
  <Box as={BOX_TAG.HEADER} className={styles['school-home__header']}>
    <Box as={BOX_TAG.DIV}>
      <Box as={BOX_TAG.DIV} className={styles['school-home__breadcrumb']}>
        <School size={14} /> RDR SCHOOL <ChevronRight size={13} /> OVERVIEW
      </Box>
      <Box as={BOX_TAG.H1}>
        Good {date.greeting}
        <Box as={BOX_TAG.SPAN} className={styles['school-home__spark']}>
          ✳
        </Box>
      </Box>
      <Box as={BOX_TAG.P}>Your school day, in one clear view.</Box>
    </Box>
    <Box as={BOX_TAG.DIV} className={styles['school-home__date']}>
      <CalendarDays size={17} />
      <Box as={BOX_TAG.SPAN}>
        <Box as={BOX_TAG.STRONG}>{date.dayName}</Box>
        <Box as={BOX_TAG.SMALL}>{date.dateLabel}</Box>
      </Box>
    </Box>
  </Box>
);
