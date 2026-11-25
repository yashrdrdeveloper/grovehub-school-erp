import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { Info, Search } from '@Components/icons';

import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../../../Button';
import { DebouncedInput } from '../../../DebouncedInput';
import { Icon } from '../../../Icon';
import { Popover } from '../../../Popover';

import type { TableSearchViewProps } from './types/tableSearch.props';
import styles from './TableSearch.module.css';

export const TableSearch = ({ value, placeholder, debounceMs, fields, onChange }: TableSearchViewProps) => {
  const [open, setOpen] = useState(false);
  return (
    <Box className={styles.root}>
      <Box className={styles.input}>
        <Box className={styles.icon} aria-hidden="true"><Icon icon={Search} size={14} /></Box>
        <DebouncedInput type="search" label={__('Search')} value={value} placeholder={placeholder} debounceMs={debounceMs} onChange={onChange} />
      </Box>
      {!!fields?.length && (
        <Popover open={open} onOpenChange={setOpen} title={__('Searches these fields')} width={248} trigger={
          <Button className={styles.info} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Searchable fields')} aria-expanded={open} onClick={() => setOpen(true)}>
            <Icon icon={Info} size={14} />
          </Button>
        }>
          <Box className={styles.fields}>
            {_map(fields, (field) => <Box key={field} className={styles.field}>{field}</Box>)}
          </Box>
        </Popover>
      )}
    </Box>
  );
};
