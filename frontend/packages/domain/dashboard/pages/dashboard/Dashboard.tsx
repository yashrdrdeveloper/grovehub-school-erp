import { DashboardHeader } from './components/dashboardHeader';
import { DashboardNotice } from './components/dashboardNotice';
import { DashboardOverview } from './components/dashboardOverview';
import { DashboardMain } from './components/dashboardMain';
import { DashboardAside } from './components/dashboardAside';
import type { DashboardComponent } from './types';

/** Pure dashboard composition. All data is explicitly marked as a preview. */
export const Dashboard: DashboardComponent = props => (
  <div className="school-home">
    <DashboardHeader date={props.date}/>
    <DashboardNotice/>
    <DashboardOverview/>
    <div className="school-home__layout">
      <DashboardMain/>
      <DashboardAside {...props}/>
    </div>
  </div>
);
