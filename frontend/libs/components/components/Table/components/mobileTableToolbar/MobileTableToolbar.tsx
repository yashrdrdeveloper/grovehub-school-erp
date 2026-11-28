import { useState } from 'react';
import { _map } from '@Base/utilities/map';
import { Settings2 } from '@Components/icons';
import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';
import { Filters } from '../../../Filters';
import { Sort } from '../../../Sort';
import { ColumnManager } from '../../../ColumnManager';
import { Sheet } from '../../../Sheet';
import { MobileTableOptions } from './components/mobileTableOptions';
import { TableSearch } from '../tableSearch';
import { MobileTablePanel } from './constants/mobileTableToolbar.panel';
import type { MobileTableToolbarProps } from './types/mobileTableToolbar.props';
import styles from './MobileTableToolbar.module.css';

export const MobileTableToolbar = <T,>(props: MobileTableToolbarProps<T>) => {
  const {
    searchProps,
    filterProps,
    sortProps,
    paginationProps,
    options,
    controls,
    actions,
    mobileRowProps,
    mobileInfiniteProps,
    maxPinnedColumns,
  } = props;
  const [panel, setPanel] = useState<MobileTablePanel | null>(null);
  const close = () => setPanel(null);
  const getBatchLabel = () => {
    if (mobileInfiniteProps) return __('Records per load');
    return __('Rows per page');
  };
  return (
    <Box className={styles.root}>
      {searchProps && (
        <TableSearch
          {...searchProps}
          placeholder={searchProps.mobilePlaceholder ?? searchProps.placeholder}
          fields={undefined}
          onChange={(value) => actions?.SEARCH?.(value)}
        />
      )}
      <ButtonWithIcon
        className={styles.trigger}
        icon={Settings2}
        iconSize={18}
        size={BUTTON_SIZE.SMALL}
        variant={BUTTON_VARIANT.OUTLINE}
        aria-label={__('Table options')}
        aria-expanded={panel !== null}
        onClick={() => setPanel(MobileTablePanel.OPTIONS)}
      />
      <MobileTableOptions
        {...props}
        panel={panel}
        setPanel={setPanel}
        close={close}
        batchLabel={getBatchLabel()}
      />
      {panel === MobileTablePanel.FILTER && filterProps && (
        <Filters
          {...filterProps}
          initialOpen
          onDismiss={close}
          onChange={(value) => actions?.FILTER?.(value)}
        />
      )}
      {panel === MobileTablePanel.SORT && sortProps && (
        <Sort
          {...sortProps}
          initialOpen
          onDismiss={close}
          onChange={(value) => actions?.SORT?.(value)}
        />
      )}
      {panel === MobileTablePanel.SETTINGS && (
        <ColumnManager
          initialOpen
          onDismiss={close}
          mobileList={Boolean(mobileRowProps)}
          options={options}
          density={controls.preferences.density}
          onToggle={controls.toggle}
          onDensityChange={controls.setRowDensity}
          onReset={controls.reset}
          onReorder={controls.reorder}
          onPin={controls.pin}
          maxPinnedColumns={maxPinnedColumns}
        />
      )}
      {panel === MobileTablePanel.BATCH && paginationProps && (
        <Sheet
          open
          onOpenChange={(open) => {
            if (!open) close();
          }}
          title={getBatchLabel()}
        >
          <Box className={styles.options}>
            {_map(paginationProps.pageSizeOptions ?? [10, 25, 50, 100], (size) => (
              <Button
                key={size}
                variant={BUTTON_VARIANT.OUTLINE}
                aria-pressed={size === paginationProps.pageSize}
                onClick={() => {
                  actions?.PAGINATION?.({ type: 'PAGE_SIZE_CHANGE', pageSize: size });
                  close();
                }}
              >
                {size}
              </Button>
            ))}
          </Box>
        </Sheet>
      )}
    </Box>
  );
};
