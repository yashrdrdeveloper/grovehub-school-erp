import type { DASHBOARD_TONE } from '../dashboard.tone';

export type DashboardTone = (typeof DASHBOARD_TONE)[keyof typeof DASHBOARD_TONE];
