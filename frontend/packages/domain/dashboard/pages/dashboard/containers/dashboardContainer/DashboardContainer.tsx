import { Dashboard } from '../../Dashboard';
import { getDashboardDateParts } from '../../helpers/getDashboardDateParts';
import { useDashboardCalendar } from '../../hooks/useDashboardCalendar';
import type { DashboardContainerComponent } from './types/dashboardContainer';

/** Owns dates and interactive calendar selection; the page has no state. */
export const DashboardContainer: DashboardContainerComponent = () => {
  const calendar = useDashboardCalendar();
  const date = getDashboardDateParts(new Date());
  return <Dashboard calendar={calendar} date={date}/>;
};
