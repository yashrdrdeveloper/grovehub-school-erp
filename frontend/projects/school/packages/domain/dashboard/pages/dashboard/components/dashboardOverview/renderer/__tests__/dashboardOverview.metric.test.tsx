import { render, screen } from '@testing-library/react';

import { DASHBOARD_METRICS } from '../../../../constants/dashboard.metrics';
import { renderDashboardOverviewMetric } from '../dashboardOverview.metric';

describe('renderDashboardOverviewMetric', () => {
  it('renders an accessible unavailable metric and its description', () => {
    render(renderDashboardOverviewMetric(DASHBOARD_METRICS[0]));

    expect(screen.getByLabelText('Tasks for you not available').textContent).toContain('—');
    expect(screen.getByText('Your work queue')).toBeTruthy();
  });
});
