import { DASHBOARD_METRICS } from '../../config/dashboardPreview.config';
import type { DashboardOverviewComponent } from './types/dashboardOverview';

export const DashboardOverview: DashboardOverviewComponent = () => (
  <section className="school-home__section" aria-labelledby="home-priorities">
    <div className="school-home__section-heading">
      <div><span className="school-home__eyebrow">YOUR DAY / 01</span><h2 id="home-priorities">At a glance</h2></div>
      <span className="school-home__hint">No live data yet</span>
    </div>
    <div className="school-home__stats">
      {DASHBOARD_METRICS.map(({ name, description, icon: Icon, tone }) => (
        <div className={'school-home__stat school-home__tone--' + tone} key={name}>
          <div className="school-home__stat-label"><span>{name}</span><Icon size={18}/></div>
          <strong aria-label={name + ' not available'}>—</strong>
          <small>{description}</small>
        </div>
      ))}
    </div>
  </section>
);
