import { Inbox } from 'lucide-react';

export const EmptyState = ({ title, description }: { title: string; description: string }) => (
  <div className="empty-state">
    <Inbox size={28} strokeWidth={1.5} />
    <strong>{title}</strong>
    <p>{description}</p>
  </div>
);
