import { CalendarDays, ChevronRight, School } from 'lucide-react';

import type { DashboardHeaderComponent } from './types/dashboardHeader';

/** Presentational greeting and local date. */
export const DashboardHeader: DashboardHeaderComponent = ({ date }) => (
  <header className="school-home__header">
    <div>
      <div className="school-home__breadcrumb"><School size={14}/> RDR SCHOOL <ChevronRight size={13}/> OVERVIEW</div>
      <h1>Good {date.greeting}<span className="school-home__spark">✳</span></h1>
      <p>Your school day, in one clear view.</p>
    </div>
    <div className="school-home__date"><CalendarDays size={17}/><span><strong>{date.dayName}</strong><small>{date.dateLabel}</small></span></div>
  </header>
);
