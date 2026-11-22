import { Box } from '@Components/components/Box';
import { CircleDashed, Sparkles } from '@Components/icons';
import type { DashboardNoticeComponent } from './types/dashboardNotice';

export const DashboardNotice: DashboardNoticeComponent = () => (
  <Box as="div" className="school-home__notice" role="status">
    <Box as="div" className="school-home__notice-mark"><Sparkles size={17}/></Box>
    <Box as="div"><Box as="strong">Your workspace is taking shape</Box>
      <Box as="span">This is a preview of the dashboard layout. Your school's activity will appear as services are connected.</Box>
    </Box>
    <Box as="span" className="school-home__notice-state"><CircleDashed size={14}/> PREVIEW</Box>
  </Box>
);
