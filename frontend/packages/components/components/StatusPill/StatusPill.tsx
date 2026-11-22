import { cn } from '@Base/utilities';
import { getStatusPillTone } from './helpers/getStatusPillTone';
import type { StatusPillComponent } from './types/statusPill.props';

export const StatusPill: StatusPillComponent = ({ value }) => (
  <span className={cn('status-pill', `status-pill--${getStatusPillTone(value)}`)}>
    <span className="status-pill__dot" aria-hidden="true"/>
    {value.replaceAll('_', ' ').toLowerCase()}
  </span>
);
