import { GripVertical } from '@Components/icons';
import { Box } from '../Box';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { Tooltip } from '../Tooltip';
import { List } from '../List';
import type { DragListViewProps } from './types/dragList.view';
import styles from './DragList.module.css';

export const DragListView = <T,>({ controls, items, getItemId, renderItem, label, emptyContent, dataTestId, disabled, disabledReason, compact, canMove }: DragListViewProps<T>) => (
  <List items={items} getItemId={getItemId} label={label} emptyContent={emptyContent} dataTestId={dataTestId} renderItem={(item, index) => {
    const id = getItemId(item);
    const movable = !disabled && (!canMove || items.some((other) => other !== item && canMove(item, other)));
    const getHandleSize = () => {
      if (compact) return BUTTON_SIZE.SMALL;
      return BUTTON_SIZE.MEDIUM;
    };
    return (
      <Box className={styles.row} data-compact={compact} data-active={controls.active === id} data-target={controls.target === id && controls.active !== id}
        onDragOver={(event) => {
          if (!controls.active) return;
          event.preventDefault();
          controls.setTarget(id);
        }}
        onDrop={(event) => {
          event.preventDefault();
          if (controls.active) controls.move(controls.active, index);
          controls.finish();
        }}>
        <Tooltip label={movable && __('Drag to reorder') || disabledReason} focusable={!movable && Boolean(disabledReason)}>
        <ButtonWithIcon icon={GripVertical} size={getHandleSize()} variant={BUTTON_VARIANT.GHOST} className={styles.handle} disabled={!movable}
          draggable={movable} aria-label={__('Reorder item; use up and down arrows')}
          onDragStart={(event) => {
            event.dataTransfer.effectAllowed = 'move';
            event.dataTransfer.setData('text/plain', id);
            controls.setActive(id);
          }}
          onDragEnd={controls.finish}
          onKeyDown={(event) => {
            if (event.key === 'ArrowUp') {
              event.preventDefault();
              controls.move(id, index - 1);
            }
            if (event.key === 'ArrowDown') {
              event.preventDefault();
              controls.move(id, index + 1);
            }
          }} />
        </Tooltip>
        <Box className={styles.content}>{renderItem(item, index)}</Box>
      </Box>
    );
  }} />
);
