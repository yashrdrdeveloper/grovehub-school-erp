import { Inbox } from '@Components/icons';

import { Box } from '../Box';
import { Icon } from '../Icon';
import type { EmptyStateComponent } from './types/emptyState.props';
import styles from './EmptyState.module.css';

export const EmptyState: EmptyStateComponent = ({
  title,
  description,
  icon = Inbox,
  compact = false,
  dataTestId,
}) => (
  <Box className={styles.root} data-compact={compact} dataTestId={dataTestId}>
    <Box className={styles.illustration} aria-hidden="true">
      <Box className={styles.symbol}>
        <Icon icon={icon} size={24} />
      </Box>
    </Box>
    <Box className={styles.title}>{title}</Box>
    <Box className={styles.description}>{description}</Box>
  </Box>
);
