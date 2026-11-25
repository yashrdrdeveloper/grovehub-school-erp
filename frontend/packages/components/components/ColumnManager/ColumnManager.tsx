import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { Settings, Lock, LockOpen } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Checkbox } from '../Checkbox';
import { DragList } from '../DragList';
import { Tooltip } from '../Tooltip';
import { Popover } from '../Popover';
import { Select } from '../Select';
import { TableDensity, TablePin } from '../Table/constants/table.type';

import styles from './ColumnManager.module.css';
import type { ColumnManagerProps } from './types/columnManager.props';

export const ColumnManager = ({ options, density, onToggle, onDensityChange, onReset, onReorder, onPin, maxPinnedColumns = 3, dataTestId }: ColumnManagerProps) => {
  const [open, setOpen] = useState(false);
  const isLocked = (option: ColumnManagerProps['options'][number]) => Boolean(option.pin && option.pin !== TablePin.NONE);
  const lockedCount = options.filter(isLocked).length;
  return (
    <Box dataTestId={dataTestId}>
      <Popover open={open} onOpenChange={setOpen} title={__('Table settings')} width={320} compact trigger={
        <Tooltip label={__('Table settings')}><ButtonWithIcon icon={Settings} iconSize={16} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Table settings')} aria-expanded={open} onClick={() => setOpen(true)} /></Tooltip>
      }>
        <Box className={styles.content}>
          <Select label={__('View')} value={density} clearable={false} searchable={false} options={[{value: TableDensity.COMPACT, label: __('Compact')}, {value: TableDensity.DEFAULT, label: __('Default')}, {value: TableDensity.COMFORTABLE, label: __('Comfortable')}]} onChange={(value) => { if (value) onDensityChange(value); }} />
          <Box className={styles.heading}>{__('Columns')}</Box>
          <Box className={styles.hint}>{__('Drag to reorder · {{count}}/{{max}} locked', { count: lockedCount, max: maxPinnedColumns })}</Box>
          <DragList items={options} getItemId={(option) => option.id} label={__('Column order')} disabled={!onReorder}
            canMove={(item, target) => item.reorderable !== false && target.reorderable !== false}
            onReorder={(items) => onReorder?.(_map(items, (item) => item.id))}
            renderItem={(option) => (
              <Box className={styles.option}>
                <Checkbox checked={option.visible} disabled={option.disabled} label={option.label} onChange={() => onToggle(option.id)} />
                <Box className={styles.name}>{option.label}</Box>
                {onPin && isLocked(option) && <ButtonWithIcon icon={Lock} iconSize={14} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} className={styles.locked} aria-label={__('Unlock {{column}}', { column: option.label })} aria-pressed={true} onClick={() => onPin(option.id, TablePin.NONE)} />}
                {onPin && !isLocked(option) && <ButtonWithIcon icon={LockOpen} iconSize={14} className={styles.unlocked} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} disabled={!option.visible || lockedCount >= maxPinnedColumns} aria-label={__('Lock {{column}}', { column: option.label })} aria-pressed={false} onClick={() => onPin(option.id, TablePin.LEFT)} />}
              </Box>
            )} />
          <Button variant={BUTTON_VARIANT.OUTLINE} size={BUTTON_SIZE.SMALL} onClick={onReset}>{__('Reset columns')}</Button>
        </Box>
      </Popover>
    </Box>
  );
};
