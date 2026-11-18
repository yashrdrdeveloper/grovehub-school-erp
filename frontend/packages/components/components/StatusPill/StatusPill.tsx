import { cn } from '@Base/utilities';

const tones: Record<string, string> = {
  ACTIVE: 'positive',
  PENDING_SETUP: 'notice',
  SUSPENDED: 'warning',
  DEACTIVATED: 'neutral',
};

export const StatusPill = ({ value }: { value: string }) => (
  <span className={cn('status-pill', `status-pill--${tones[value] ?? 'neutral'}`)}>
    <span className="status-pill__dot" aria-hidden="true" />
    {value.replaceAll('_', ' ').toLowerCase()}
  </span>
);
