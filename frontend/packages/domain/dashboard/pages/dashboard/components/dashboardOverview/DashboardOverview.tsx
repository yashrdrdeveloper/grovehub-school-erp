import { _map } from '@Base/utilities/map';
import { Box } from '@Components/components/Box';
import { DASHBOARD_METRICS } from '../../constants/dashboard.metrics';
import type { DashboardOverviewComponent } from './types/dashboardOverview';

export const DashboardOverview: DashboardOverviewComponent = () => (
  <Box as="section" className="school-home__section" aria-labelledby="home-priorities">
    <Box as="div" className="school-home__section-heading">
      <Box as="div"><Box as="span" className="school-home__eyebrow">YOUR DAY / 01</Box><Box as="h2" id="home-priorities">At a glance</Box></Box>
      <Box as="span" className="school-home__hint">No live data yet</Box>
    </Box>
    <Box as="div" className="school-home__stats">
      {_map(DASHBOARD_METRICS, ({ name, description, icon: Icon, tone }) => (
        <Box as="div" className={'school-home__stat school-home__tone--' + tone} key={name}>
          <Box as="div" className="school-home__stat-label"><Box as="span">{name}</Box><Icon size={18}/></Box>
          <Box as="strong" aria-label={name + ' not available'}>—</Box>
          <Box as="small">{description}</Box>
        </Box>
      ))}
    </Box>
  </Box>
);
