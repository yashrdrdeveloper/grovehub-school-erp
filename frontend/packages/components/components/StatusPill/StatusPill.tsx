import { cn } from '@Base/utilities/cn';

import { Box } from '../Box';
import { getStatusPillDisplayLabel } from './helpers/statusPill.getDisplayLabel';
import { getStatusPillTone } from './helpers/statusPill.getTone';

import type { StatusPillComponent } from './types/statusPill.props';

/** Stateless status presentation with a shared Box and Lodash-derived label. */
export const StatusPill: StatusPillComponent = ({ value, dataTestId }) => (
  <Box
    as="span"
    dataTestId={dataTestId}
    className={cn('status-pill', `status-pill--${getStatusPillTone(value)}`)}
  >
    <Box as="span" className="status-pill__dot" aria-hidden="true" />
    {getStatusPillDisplayLabel(value)}
  </Box>
);
