import type { ReactElement } from 'react';
import type { DASHBOARD_METRICS } from '../../../../constants/dashboard.metrics';

export type DashboardMetric = (typeof DASHBOARD_METRICS)[number];
export type RenderDashboardOverviewMetric = (metric: DashboardMetric) => ReactElement;
