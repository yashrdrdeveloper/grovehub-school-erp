import { cn } from '@Base/utilities/cn';
import { getStatusPillTone } from './helpers/statusPill.getTone';
import type { StatusPillComponent } from './types/statusPill.props';

export const StatusPill: StatusPillComponent = ({ value, dataTestId }) => (
  <span data-testid={dataTestId} className={cn('status-pill', `status-pill--${getStatusPillTone(value)}`)}>
    <span className="status-pill__dot" aria-hidden="true"/>
    {value.replaceAll('_', ' ').toLowerCase()}
  </span>
);
