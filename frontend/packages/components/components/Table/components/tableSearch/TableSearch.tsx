import { useState } from 'react';
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
  const searchDescription = __('Search is based on {{fields}}.', { fields: fields?.join(', ') ?? '' });
  return (
    <Box className={styles.root}>
      <Box className={styles.input}>
        <Box className={styles.icon} aria-hidden="true"><Icon icon={Search} size={14} /></Box>
        <DebouncedInput type="search" label={__('Search')} value={value} placeholder={placeholder} debounceMs={debounceMs} onChange={onChange} />
      </Box>
      {!!fields?.length && (
        <Popover open={open} onOpenChange={setOpen} title={__('Search information')} width={260} compact desktopDropdown desktopDropdownWidth={260} trigger={
          <Button className={styles.info} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Searchable fields')} title={searchDescription} aria-expanded={open} onClick={() => setOpen(true)}>
            <Icon icon={Info} size={14} />
          </Button>
        }>
          <Box className={styles.description}>{searchDescription}</Box>
        </Popover>
      )}
    </Box>
  );
};
