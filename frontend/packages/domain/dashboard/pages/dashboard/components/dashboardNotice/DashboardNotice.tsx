import { CircleDashed, Sparkles } from '@Components/icons';
import type { DashboardNoticeComponent } from './types/dashboardNotice';

export const DashboardNotice: DashboardNoticeComponent = () => (
  <div className="school-home__notice" role="status">
    <div className="school-home__notice-mark"><Sparkles size={17}/></div>
    <div><strong>Your workspace is taking shape</strong>
      <span>This is a preview of the dashboard layout. Your school's activity will appear as services are connected.</span>
    </div>
    <span className="school-home__notice-state"><CircleDashed size={14}/> PREVIEW</span>
  </div>
);
