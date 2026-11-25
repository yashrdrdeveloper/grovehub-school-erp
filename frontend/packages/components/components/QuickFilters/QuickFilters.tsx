import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';

import styles from './QuickFilters.module.css';
import type { QuickFiltersProps } from './types/quickFilters.props';

export const QuickFilters = ({ options, value, multiple, onChange, dataTestId }: QuickFiltersProps) => {
  const toggle = (id: string) => {
    if (!multiple) {
      onChange([id]);
      return;
    }
    if (value.includes(id)) {
      onChange(_filter(value, (entry) => entry !== id));
      return;
    }
    onChange([...value, id]);
  };
  return (
    <Box className={styles.controls} dataTestId={dataTestId} role="group" aria-label={__('Quick filters')}>
      {_map(options, (option) => {
        const variant = () => {
          if (value.includes(option.id)) return BUTTON_VARIANT.SECONDARY;
          return BUTTON_VARIANT.GHOST;
        };
        return (
          <Button key={option.id} className={styles.filter} size={BUTTON_SIZE.SMALL} variant={variant()} aria-pressed={value.includes(option.id)} disabled={option.disabled} onClick={() => toggle(option.id)}>
            {option.label}
            {option.count !== undefined && ` (${option.count})`}
          </Button>
        );
      })}
    </Box>
  );
};
