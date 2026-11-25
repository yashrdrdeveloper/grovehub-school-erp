import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { Settings2 } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Sheet } from '../Sheet';
import { Select } from '../Select';
import { TableDensity } from '../Table/constants/table.type';

import styles from './ColumnManager.module.css';
import type { ColumnManagerProps } from './types/columnManager.props';

export const ColumnManager = ({ options, density, onToggle, onDensityChange, onReset, dataTestId }: ColumnManagerProps) => {
  const [open, setOpen] = useState(false);
  return (
    <Box dataTestId={dataTestId}>
      <ButtonWithIcon icon={Settings2} variant={BUTTON_VARIANT.OUTLINE} onClick={() => setOpen(true)}>{__('Columns')}</ButtonWithIcon>
      <Sheet open={open} onOpenChange={setOpen} title={__('Columns')}>
        <Box className={styles.content}>
          {_map(options, (option) => (
            <Button key={option.id} variant={BUTTON_VARIANT.OUTLINE} aria-pressed={option.visible} disabled={option.disabled} onClick={() => onToggle(option.id)}>
              {option.visible && '✓ '}{option.label}
            </Button>
          ))}
          <Select label={__('Density')} value={density} clearable={false} searchable={false} options={[{value: TableDensity.COMPACT, label: __('Compact')}, {value: TableDensity.DEFAULT, label: __('Default')}, {value: TableDensity.COMFORTABLE, label: __('Comfortable')}]} onChange={(value) => {
            if (value) onDensityChange(value);
          }} />
          <Button variant={BUTTON_VARIANT.OUTLINE} onClick={onReset}>{__('Reset columns')}</Button>
        </Box>
      </Sheet>
    </Box>
  );
};
