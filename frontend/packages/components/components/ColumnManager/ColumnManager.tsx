import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { Settings2 } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Popover } from '../Popover';
import { Select } from '../Select';
import { TableDensity, TablePin } from '../Table/constants/table.type';

import styles from './ColumnManager.module.css';
import type { ColumnManagerProps } from './types/columnManager.props';

export const ColumnManager = ({ options, density, onToggle, onDensityChange, onReset, onMove, onResize, onPin, dataTestId }: ColumnManagerProps) => {
  const [open, setOpen] = useState(false);
  return (
    <Box dataTestId={dataTestId}>
      <Popover open={open} onOpenChange={setOpen} title={__('Columns')} width={340} trigger={
      <ButtonWithIcon icon={Settings2} variant={BUTTON_VARIANT.OUTLINE} onClick={() => setOpen(true)}>{__('Columns')}</ButtonWithIcon>
      }>
        <Box className={styles.content}>
          {_map(options, (option, index) => (
            <Box className={styles.option} key={option.id}>
              <Button variant={BUTTON_VARIANT.OUTLINE} aria-pressed={option.visible} disabled={option.disabled} onClick={() => onToggle(option.id)}>{option.visible && '✓ '}{option.label}</Button>
              <Box className={styles.actions}>
                {onMove && <Button variant={BUTTON_VARIANT.GHOST} disabled={option.reorderable === false || index === 0} aria-label={__('Move column up')} onClick={() => onMove(option.id, -1)}>↑</Button>}
                {onMove && <Button variant={BUTTON_VARIANT.GHOST} disabled={option.reorderable === false || index === options.length - 1} aria-label={__('Move column down')} onClick={() => onMove(option.id, 1)}>↓</Button>}
                {onResize && option.resizable !== false && <Button variant={BUTTON_VARIANT.GHOST} aria-label={__('Decrease width')} onClick={() => onResize(option.id, (option.width ?? 180) - 20)}>−</Button>}
                {onResize && option.resizable !== false && <Button variant={BUTTON_VARIANT.GHOST} aria-label={__('Increase width')} onClick={() => onResize(option.id, (option.width ?? 180) + 20)}>+</Button>}
                {onResize && <Box>{option.width ?? 180}px</Box>}
              </Box>
              {onPin && <Select label={__('Pin column')} value={option.pin ?? TablePin.NONE} clearable={false} searchable={false} options={[{value: TablePin.NONE, label: __('Not pinned')}, {value: TablePin.LEFT, label: __('Pin left')}, {value: TablePin.RIGHT, label: __('Pin right')}]} onChange={(value) => { if (value) onPin(option.id, value); }} />}
            </Box>
          ))}
          <Select label={__('Density')} value={density} clearable={false} searchable={false} options={[{value: TableDensity.COMPACT, label: __('Compact')}, {value: TableDensity.DEFAULT, label: __('Default')}, {value: TableDensity.COMFORTABLE, label: __('Comfortable')}]} onChange={(value) => { if (value) onDensityChange(value); }} />
          <Button variant={BUTTON_VARIANT.OUTLINE} onClick={onReset}>{__('Reset columns')}</Button>
        </Box>
      </Popover>
    </Box>
  );
};
