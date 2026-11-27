import { useState } from 'react';
import { Pencil } from '@Components/icons';

import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';

import type { TableCellProps } from './types/tableCell.props';
import styles from './TableCell.module.css';

export const TableCell = <T,>({ row, column, actions }: TableCellProps<T>) => {
  const [editing, setEditing] = useState(false);
  const render = () => {
    if (column.renderer) return column.renderer(row);
    return String(column.accessor(row) ?? '—');
  };
  const content = render();
  const getTitle = () => {
    if (typeof content === 'string') return content;
    return undefined;
  };
  if (editing && column.editRenderer && actions?.CELL_EDIT) {
    return (
      <Box className={styles.cell} onClick={(event) => event.stopPropagation()}>
        {column.editRenderer(row, (value) => {
          actions.CELL_EDIT?.({ row, columnId: column.id, value });
          setEditing(false);
        })}
        <Button variant={BUTTON_VARIANT.OUTLINE} onClick={() => setEditing(false)}>
          {__('Cancel')}
        </Button>
      </Box>
    );
  }
  return (
    <Box className={styles.cell}>
      {actions?.CELL_CLICK && (
        <Button
          variant={BUTTON_VARIANT.GHOST}
          onClick={(event) => {
            event.stopPropagation();
            actions.CELL_CLICK?.({ row, columnId: column.id });
          }}
        >
          {content}
        </Button>
      )}
      {!actions?.CELL_CLICK && (
        <Box className={styles.value} title={getTitle()}>
          {content}
        </Box>
      )}
      {column.editRenderer && actions?.CELL_EDIT && (
        <ButtonWithIcon
          icon={Pencil}
          variant={BUTTON_VARIANT.GHOST}
          aria-label={__('Edit cell')}
          onClick={(event) => {
            event.stopPropagation();
            setEditing(true);
          }}
        />
      )}
    </Box>
  );
};
