import type { DashboardDateParts } from '../../../../../types';
export interface DashboardCalendarCopy {
  readonly calendarDay: string;
  readonly calendarHeading: string;
  readonly calendarDetail: string;
  readonly calendarEmptyTitle: string;
}
export type GetDashboardCalendarCopy = (date: DashboardDateParts, today: boolean) => DashboardCalendarCopy;
