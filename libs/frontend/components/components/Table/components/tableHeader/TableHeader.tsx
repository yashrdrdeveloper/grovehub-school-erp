import { _map } from '@Base/utilities/map';

import { Box, BOX_TAG } from '../../../Box';
import { Button } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';

import styles from './TableHeader.module.css';
import type { TableHeaderViewProps } from './types/tableHeader.props';

export const TableHeader = ({
  title,
  description,
  leftConfig,
  rightConfig,
  onAction,
}: TableHeaderViewProps) => {
  const actions = (items: TableHeaderViewProps['rightConfig']) =>
    _map(items ?? [], (item) => {
      if (item.icon)
        return (
          <ButtonWithIcon
            key={item.id}
            icon={item.icon}
            disabled={item.disabled}
            aria-label={item.label}
            onClick={() => onAction?.(item.id)}
          >
            <Box as={BOX_TAG.SPAN} className={styles.iconActionLabel}>
              {item.label}
            </Box>
          </ButtonWithIcon>
        );
      return (
        <Button key={item.id} disabled={item.disabled} onClick={() => onAction?.(item.id)}>
          {item.label}
        </Button>
      );
    });
  return (
    <Box className={styles.header}>
      <Box>
        {title && (
          <Box as={BOX_TAG.H1} className={styles.title}>
            {title}
          </Box>
        )}
        {description && (
          <Box as={BOX_TAG.P} className={styles.description}>
            {description}
          </Box>
        )}
        {actions(leftConfig)}
      </Box>
      <Box className={styles.actions}>{actions(rightConfig)}</Box>
    </Box>
  );
};
