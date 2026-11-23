import { getDashboardOverviewUnavailableLabel } from '../dashboardOverview.getUnavailableLabel';

describe('getDashboardOverviewUnavailableLabel', () => {
  it('interpolates metric accessibility copy', () => {
    expect(getDashboardOverviewUnavailableLabel('Tasks for you')).toBe('Tasks for you not available');
  });
});
