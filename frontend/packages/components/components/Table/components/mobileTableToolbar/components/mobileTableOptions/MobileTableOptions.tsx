import { Filter, ArrowUpDown, Settings, RefreshCw, Layers2, Download } from '@Components/icons';
import { Box } from '../../../../../Box';
import { BUTTON_VARIANT } from '../../../../../Button';
import { ButtonWithIcon } from '../../../../../ButtonWithIcon';
import { Sheet } from '../../../../../Sheet';

import { MobileTablePanel } from '../../constants/mobileTableToolbar.panel';
import styles from '../../MobileTableToolbar.module.css';
import type { MobileTableOptionsProps } from './types/mobileTableOptions.props';

export const MobileTableOptions = <T,>({
  filterProps,
  sortProps,
  paginationProps,
  actions,
  refreshing,
  panel,
  setPanel,
  close,
  batchLabel,
}: MobileTableOptionsProps<T>) => {
  const getFilterLabel = () => {
    if (filterProps?.value.length)
      return __('Filters ({{count}})', { count: filterProps.value.length });
    return __('Filters');
  };
  const getSortLabel = () => {
    if (sortProps?.value.length) return __('Sort ({{count}})', { count: sortProps.value.length });
    return __('Sort');
  };
  return (
    <Sheet
      open={panel === MobileTablePanel.OPTIONS}
      onOpenChange={(open) => {
        if (!open) close();
      }}
      title={__('Table options')}
    >
      <Box className={styles.options}>
        {filterProps && (
          <ButtonWithIcon
            icon={Filter}
            variant={BUTTON_VARIANT.GHOST}
            onClick={() => setPanel(MobileTablePanel.FILTER)}
          >
            {getFilterLabel()}
          </ButtonWithIcon>
        )}
        {sortProps && (
          <ButtonWithIcon
            icon={ArrowUpDown}
            variant={BUTTON_VARIANT.GHOST}
            onClick={() => setPanel(MobileTablePanel.SORT)}
          >
            {getSortLabel()}
          </ButtonWithIcon>
        )}
        <ButtonWithIcon
          icon={Settings}
          variant={BUTTON_VARIANT.GHOST}
          onClick={() => setPanel(MobileTablePanel.SETTINGS)}
        >
          {__('Table settings')}
        </ButtonWithIcon>
        {paginationProps && (
          <ButtonWithIcon
            icon={Layers2}
            variant={BUTTON_VARIANT.GHOST}
            onClick={() => setPanel(MobileTablePanel.BATCH)}
          >
            {batchLabel} · {paginationProps.pageSize}
          </ButtonWithIcon>
        )}
        {actions?.REFRESH && (
          <ButtonWithIcon
            icon={RefreshCw}
            variant={BUTTON_VARIANT.GHOST}
            disabled={refreshing}
            onClick={() => {
              close();
              actions.REFRESH?.();
            }}
          >
            {__('Refresh')}
          </ButtonWithIcon>
        )}
        {actions?.EXPORT && (
          <ButtonWithIcon
            icon={Download}
            variant={BUTTON_VARIANT.GHOST}
            onClick={() => {
              close();
              actions.EXPORT?.();
            }}
          >
            {__('Export')}
          </ButtonWithIcon>
        )}
      </Box>
    </Sheet>
  );
};
