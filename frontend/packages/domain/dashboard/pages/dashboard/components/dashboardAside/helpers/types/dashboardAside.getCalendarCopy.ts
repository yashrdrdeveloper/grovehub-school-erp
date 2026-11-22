import type { DashboardAsideProps } from '../../types/dashboardAside';
export interface DashboardCalendarCopy {
  readonly calendarDay: string;
  readonly calendarHeading: string;
  readonly calendarDetail: string;
  readonly calendarEmptyTitle: string;
}
export type GetDashboardCalendarCopy = (date: DashboardAsideProps['date'], today: boolean) => DashboardCalendarCopy;
