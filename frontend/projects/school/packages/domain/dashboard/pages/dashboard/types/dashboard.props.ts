import type { ReactElement } from 'react';

export type DashboardCalendarView = 'today' | 'upcoming';

export interface DashboardDateParts {
  readonly greeting: 'morning' | 'afternoon' | 'evening';
  readonly dayName: string;
  readonly dateLabel: string;
  readonly dayOfMonth: number;
}
export interface DashboardCalendarControls {
  readonly view: DashboardCalendarView;
  readonly showToday: () => void;
  readonly showUpcoming: () => void;
}
export interface DashboardProps {
  readonly date: DashboardDateParts;
  readonly calendar: DashboardCalendarControls;
}
export type DashboardComponent = (props: DashboardProps) => ReactElement;
