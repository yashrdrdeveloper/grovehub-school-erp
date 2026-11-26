import { RefreshCw } from '@Components/icons';
import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';
import { ColumnManager } from '../../../ColumnManager';
import { Filters } from '../../../Filters';
import { Sort } from '../../../Sort';
import { Tooltip } from '../../../Tooltip';

import { TableSearch } from '../tableSearch';
import { TableOptions } from '../tableOptions';
import styles from '../../Table.module.css';
import type { DesktopTableToolbarProps } from './types/desktopTableToolbar.props';

export const DesktopTableToolbar = <T,>(props: DesktopTableToolbarProps<T>) => {
  const { actions, controls, options, searchProps, filterProps, sortProps } = props;
  const getRefreshLabel = () => {
    if (props.refreshing) return __('Refreshing…');
    return __('Refresh');
  };
  return (
    <Box className={styles.toolbar}>
      <Box className={styles.toolbarLeft}>
        {filterProps && (
          <Filters {...filterProps} iconOnly onChange={(value) => actions?.FILTER?.(value)} />
        )}
        {sortProps && <Sort {...sortProps} iconOnly onChange={(value) => actions?.SORT?.(value)} />}
      </Box>
      <Box className={styles.toolbarRight}>
        {searchProps && (
          <TableSearch {...searchProps} onChange={(value) => actions?.SEARCH?.(value)} />
        )}
        {actions?.REFRESH && (
          <Tooltip label={getRefreshLabel()} focusable={Boolean(props.refreshing)}>
            <ButtonWithIcon
              icon={RefreshCw}
              iconSize={16}
              size={BUTTON_SIZE.SMALL}
              variant={BUTTON_VARIANT.OUTLINE}
              aria-label={__('Refresh')}
              onClick={actions.REFRESH}
              disabled={props.refreshing}
            />
          </Tooltip>
        )}
        <ColumnManager
          options={options}
          density={controls.preferences.density}
          onToggle={controls.toggle}
          onDensityChange={controls.setRowDensity}
          onReset={controls.reset}
          onReorder={controls.reorder}
          onPin={controls.pin}
          maxPinnedColumns={props.maxPinnedColumns}
        />
        {(actions?.EXPORT || actions?.FULLSCREEN_CHANGE) && (
          <TableOptions>
            {actions?.EXPORT && (
              <Button variant={BUTTON_VARIANT.GHOST} onClick={actions.EXPORT}>
                {__('Export')}
              </Button>
            )}
            {actions?.FULLSCREEN_CHANGE && (
              <Button
                variant={BUTTON_VARIANT.GHOST}
                aria-pressed={controls.fullscreen}
                onClick={() => controls.changeFullscreen(!controls.fullscreen)}
              >
                {!controls.fullscreen && __('Fullscreen')}
                {controls.fullscreen && __('Exit fullscreen')}
              </Button>
            )}
          </TableOptions>
        )}
      </Box>
    </Box>
  );
};
