import type { ReactElement } from 'react';
import type { DASHBOARD_WORKSPACES } from '../../../../constants/dashboard.workspaces';

export type DashboardWorkspace = (typeof DASHBOARD_WORKSPACES)[number];
export type RenderDashboardMainWorkspace = (workspace: DashboardWorkspace) => ReactElement;
