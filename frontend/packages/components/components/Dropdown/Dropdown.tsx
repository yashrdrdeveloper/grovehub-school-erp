import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { ChevronDown, ChevronRight, Check } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../ButtonWithIcon';
import { Icon } from '../Icon';
import { Popover } from '../Popover';

import { DropdownItemType } from './constants/dropdown.type';
import styles from './Dropdown.module.css';
import type { DropdownItem, DropdownProps } from './types/dropdown.props';

export const Dropdown = ({ label, items, onSelect, disabled, dataTestId }: DropdownProps) => {
  const [open, setOpen] = useState(false);
  const [path, setPath] = useState<readonly string[]>([]);
  const current = path.reduce((list, id) => list.find((item) => item.id === id)?.children ?? [], items);
  const onOpenChange = (next: boolean) => {
    if (!next) setPath([]);
    setOpen(next);
  };
  const select = (item: DropdownItem) => {
    if (item.type === DropdownItemType.SUBMENU) {
      setPath([...path, item.id]);
      return;
    }
    onSelect(item.id);
    if (item.type !== DropdownItemType.CHECKBOX && item.type !== DropdownItemType.RADIO) onOpenChange(false);
  };
  return (
    <Box dataTestId={dataTestId}>
      <Popover open={open} onOpenChange={onOpenChange} title={label} width={240} trigger={<ButtonWithIcon icon={ChevronDown} iconPosition={BUTTON_ICON_POSITION.RIGHT} variant={BUTTON_VARIANT.OUTLINE} disabled={disabled} aria-haspopup="dialog" aria-expanded={open} onClick={() => onOpenChange(true)}>{label}</ButtonWithIcon>}>
        <Box className={styles.items}>
          {!!path.length && <Button variant={BUTTON_VARIANT.OUTLINE} onClick={() => setPath(path.slice(0, -1))}>{__('Back')}</Button>}
          {_map(current, (item) => {
            if (item.type === DropdownItemType.SEPARATOR) return <Box key={item.id} role="separator" className={styles.separator} />;
            if (item.type === DropdownItemType.LABEL) return <Box key={item.id} className={styles.label}>{item.label}</Box>;
            const variant = () => {
              if (item.destructive) return BUTTON_VARIANT.DANGER;
              return BUTTON_VARIANT.GHOST;
            };
            const role = () => {
              if (item.type === DropdownItemType.CHECKBOX) return 'checkbox';
              if (item.type === DropdownItemType.RADIO) return 'radio';
              return undefined;
            };
            return (
              <Button key={item.id} variant={variant()} disabled={item.disabled} role={role()} aria-checked={item.checked} onClick={() => select(item)}>
                {item.checked && <Icon icon={Check} />}
                {item.icon && <Icon icon={item.icon} />}
                <Box className={styles.itemLabel}>{item.label}</Box>
                {item.type === DropdownItemType.SUBMENU && <Icon icon={ChevronRight} />}
              </Button>
            );
          })}
        </Box>
      </Popover>
    </Box>
  );
};
