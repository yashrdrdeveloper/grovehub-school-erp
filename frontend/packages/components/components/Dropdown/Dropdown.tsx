import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { ChevronDown } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../ButtonWithIcon';
import { Icon } from '../Icon';
import { Sheet } from '../Sheet';

import styles from './Dropdown.module.css';
import type { DropdownProps } from './types/dropdown.props';

export const Dropdown = ({ label, items, onSelect, disabled, dataTestId }: DropdownProps) => {
  const [open, setOpen] = useState(false);
  return (
    <>
      <ButtonWithIcon icon={ChevronDown} iconPosition={BUTTON_ICON_POSITION.RIGHT} variant={BUTTON_VARIANT.OUTLINE} disabled={disabled} dataTestId={dataTestId} aria-haspopup="dialog" aria-expanded={open} onClick={() => setOpen(true)}>{label}</ButtonWithIcon>
      <Sheet open={open} onOpenChange={setOpen} title={label}>
        <Box className={styles.items}>
          {_map(items, (item) => {
            const getVariant = () => {
              if (item.destructive) return BUTTON_VARIANT.DANGER;
              return BUTTON_VARIANT.GHOST;
            };
            return (
              <Button key={item.id} variant={getVariant()} disabled={item.disabled} onClick={() => {
                setOpen(false);
                onSelect(item.id);
              }}>
                {item.icon && <Icon icon={item.icon} />}
                {item.label}
              </Button>
            );
          })}
        </Box>
      </Sheet>
    </>
  );
};
