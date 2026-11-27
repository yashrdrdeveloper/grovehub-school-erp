import { render, screen } from '@testing-library/react';

import { DASHBOARD_WORKSPACES } from '../../../../constants/dashboard.workspaces';
import { renderDashboardMainWorkspace } from '../dashboardMain.workspace';

describe('renderDashboardMainWorkspace', () => {
  it('renders the workspace title and preview content', () => {
    render(renderDashboardMainWorkspace(DASHBOARD_WORKSPACES[0]));

    expect(screen.getByRole('heading', { name: 'Admissions & dues' })).toBeTruthy();
    expect(screen.getByText('Applications, payments and follow-ups')).toBeTruthy();
    expect(screen.getByText('Ready for future service data')).toBeTruthy();
  });
});
