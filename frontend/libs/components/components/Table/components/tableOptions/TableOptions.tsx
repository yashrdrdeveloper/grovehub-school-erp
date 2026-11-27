import { useState } from 'react';
import { Ellipsis } from '@Components/icons';

import { Box } from '../../../Box';
import { ButtonWithIcon } from '../../../ButtonWithIcon';
import { BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { Popover } from '../../../Popover';
import { Tooltip } from '../../../Tooltip';

import type { TableOptionsProps } from './types/tableOptions.props';
import styles from './TableOptions.module.css';

export const TableOptions = ({ children }: TableOptionsProps) => {
  const [open, setOpen] = useState(false);
  return (
    <Popover
      open={open}
      onOpenChange={setOpen}
      title={__('Table options')}
      width={220}
      trigger={
        <Tooltip label={__('Table options')}>
          <ButtonWithIcon
            icon={Ellipsis}
            iconSize={16}
            variant={BUTTON_VARIANT.OUTLINE}
            size={BUTTON_SIZE.SMALL}
            aria-label={__('Table options')}
            aria-expanded={open}
            onClick={() => setOpen(true)}
          />
        </Tooltip>
      }
    >
      <Box className={styles.options}>{children}</Box>
    </Popover>
  );
};
