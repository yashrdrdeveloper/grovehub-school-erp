import { Inbox } from 'lucide-react';
import type { EmptyStateComponent } from './types/emptyState.props';

export const EmptyState: EmptyStateComponent = ({ title, description, dataTestId }) => (
  <div className="empty-state" data-testid={dataTestId}>
    <Inbox size={28} strokeWidth={1.5}/>
    <strong>{title}</strong>
    <p>{description}</p>
  </div>
);
