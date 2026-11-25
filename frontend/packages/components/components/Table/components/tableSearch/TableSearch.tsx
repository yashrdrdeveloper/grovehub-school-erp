import { Info, Search } from '@Components/icons';

import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../../../Button';
import { DebouncedInput } from '../../../DebouncedInput';
import { Icon } from '../../../Icon';
import { Tooltip } from '../../../Tooltip';

import type { TableSearchViewProps } from './types/tableSearch.props';
import styles from './TableSearch.module.css';

export const TableSearch = ({ value, placeholder, debounceMs, fields, onChange }: TableSearchViewProps) => {
  const searchDescription = __('Search is based on {{fields}}.', { fields: fields?.join(', ') ?? '' });
  return (
    <Box className={styles.root}>
      <Box className={styles.input}>
        <Box className={styles.icon} aria-hidden="true"><Icon icon={Search} size={14} /></Box>
        <DebouncedInput type="search" label={__('Search')} value={value} placeholder={placeholder} debounceMs={debounceMs} onChange={onChange} />
      </Box>
      {!!fields?.length && (
        <Tooltip label={searchDescription}>
          <Button type="button" className={styles.info} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={searchDescription}>
            <Icon icon={Info} size={14} />
          </Button>
        </Tooltip>
      )}
    </Box>
  );
};
