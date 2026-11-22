import type { ReactElement } from 'react';
import type { DashboardDateParts } from '../../../types';

export interface DashboardHeaderProps { readonly date: DashboardDateParts; }
export type DashboardHeaderComponent = (props: DashboardHeaderProps) => ReactElement;
