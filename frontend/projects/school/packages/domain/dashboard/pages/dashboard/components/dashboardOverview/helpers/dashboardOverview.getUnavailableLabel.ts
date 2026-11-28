import type { GetDashboardOverviewUnavailableLabel } from './types/dashboardOverview.getUnavailableLabel';

/** The translation source is literal; the metric name is interpolated separately. */
export const getDashboardOverviewUnavailableLabel: GetDashboardOverviewUnavailableLabel = (name) =>
  __('{{name}} not available', { name });
