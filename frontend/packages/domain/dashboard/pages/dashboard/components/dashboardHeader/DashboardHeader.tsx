import { Box } from '@Components/components/Box';
import { CalendarDays, ChevronRight, School } from '@Components/icons';

import type { DashboardHeaderComponent } from './types/dashboardHeader';

/** Presentational greeting and local date. */
export const DashboardHeader: DashboardHeaderComponent = ({ date }) => (
  <Box as="header" className="school-home__header">
    <Box as="div">
      <Box as="div" className="school-home__breadcrumb"><School size={14}/> RDR SCHOOL <ChevronRight size={13}/> OVERVIEW</Box>
      <Box as="h1">Good {date.greeting}<Box as="span" className="school-home__spark">✳</Box></Box>
      <Box as="p">Your school day, in one clear view.</Box>
    </Box>
    <Box as="div" className="school-home__date"><CalendarDays size={17}/><Box as="span"><Box as="strong">{date.dayName}</Box><Box as="small">{date.dateLabel}</Box></Box></Box>
  </Box>
);
