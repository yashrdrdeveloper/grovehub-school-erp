import { BOX_TAG } from '@Components/components/Box';
import { cn } from '@Base/utilities/cn';

import { Box } from '../Box';
import { getStatusPillDisplayLabel } from './helpers/statusPill.getDisplayLabel';
import { getStatusPillTone } from './helpers/statusPill.getTone';

import styles from './StatusPill.module.css';
import type { StatusPillComponent } from './types/statusPill.props';

/** Stateless status presentation with a shared Box and Lodash-derived label. */
export const StatusPill: StatusPillComponent = ({ value, dataTestId, plain }) => (
  <Box
    as={BOX_TAG.SPAN}
    dataTestId={dataTestId}
    data-tone={getStatusPillTone(value)}
    data-plain={plain}
    className={cn(styles.badge, 'status-pill', `status-pill--${getStatusPillTone(value)}`)}
  >
    <Box as={BOX_TAG.SPAN} className="status-pill__dot" aria-hidden="true" />
    {getStatusPillDisplayLabel(value)}
  </Box>
);
