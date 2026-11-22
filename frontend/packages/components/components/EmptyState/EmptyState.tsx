import { Inbox } from 'lucide-react';
import type { EmptyStateComponent } from './types/emptyState.props';

export const EmptyState: EmptyStateComponent = ({ title, description }) => (
  <div className="empty-state">
    <Inbox size={28} strokeWidth={1.5}/>
    <strong>{title}</strong>
    <p>{description}</p>
  </div>
);
