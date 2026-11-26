import { _filter } from '@Base/utilities/filter';
import { ChevronDown, Search } from '@Components/icons';

import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../../../ButtonWithIcon';
import { Input } from '../../../Input';
import { InputControl, INPUT_CONTROL_SIZE } from '../../../InputControl';
import { Icon } from '../../../Icon';
import { Popover } from '../../../Popover';
import { Tooltip } from '../../../Tooltip';

import { SelectionOptions } from './components/selectionOptions';
import styles from './SelectionPanel.module.css';
import type { SelectValue } from '../../types/select.props';
import type { SelectionPanelProps } from './types/selectionPanel.props';

export const SelectionPanel = <T extends SelectValue>({
  options,
  displayOptions,
  selected,
  multiple,
  open,
  search,
  onOpenChange,
  onSearch,
  onToggle,
  onClear,
  onApply,
  hasNext,
  onLoadMore,
  onRetry,
  label,
  hideLabel,
  placeholder,
  disabled,
  disabledReason,
  searchable = true,
  clearable = true,
  loading,
  error,
  dataTestId,
}: SelectionPanelProps<T>) => {
  const mobile = useMobileViewport();
  const dropdown = !mobile;
  const getDropdownWidth = () => {
    if (multiple || searchable) return 180;
    return 88;
  };
  const chosen = _filter(displayOptions ?? options, (option) => selected.includes(option.value));
  const getLabel = () => {
    if (!selected.length) return placeholder ?? __('Select an option');
    if (multiple) return __('{{count}} selected', { count: selected.length });
    return chosen[0]?.label ?? String(selected[0]);
  };

  return (
    <Box className={styles.field}>
      {!hideLabel && <Box className={styles.description}>{label}</Box>}
      <Popover
        open={open}
        onOpenChange={onOpenChange}
        title={label}
        compact
        fullWidth
        desktopDropdown
        desktopDropdownWidth={getDropdownWidth()}
        dataTestId={dataTestId && `${dataTestId}-options`}
        trigger={
          <Tooltip
            label={(disabled && disabledReason) || undefined}
            focusable={Boolean(disabled && disabledReason)}
            fullWidth
          >
            <ButtonWithIcon
              icon={ChevronDown}
              iconPosition={BUTTON_ICON_POSITION.RIGHT}
              variant={BUTTON_VARIANT.OUTLINE}
              fullWidth
              size={BUTTON_SIZE.SMALL}
              disabled={disabled}
              aria-label={label}
              aria-expanded={open}
              aria-haspopup="dialog"
              dataTestId={dataTestId}
              className={styles.trigger}
              data-placeholder={!selected.length}
              onClick={() => onOpenChange(true)}
            >
              {getLabel()}
            </ButtonWithIcon>
          </Tooltip>
        }
      >
        <Box className={styles.content} data-dropdown={dropdown} data-multiple={multiple}>
          {searchable && dropdown && (
            <Box className={styles.search}>
              <Box className={styles.searchIcon} aria-hidden="true">
                <Icon icon={Search} size={14} />
              </Box>
              <InputControl
                aria-label={__('Search options')}
                controlSize={INPUT_CONTROL_SIZE.SMALL}
                className={styles.searchInput}
                value={search}
                onChange={(event) => onSearch(event.currentTarget.value)}
                placeholder={__('Search...')}
              />
            </Box>
          )}
          {searchable && !dropdown && (
            <Input
              label={__('Search options')}
              value={search}
              onChange={(event) => onSearch(event.currentTarget.value)}
              placeholder={__('Search...')}
            />
          )}
          {error && <Box role="alert">{error}</Box>}
          {error && onRetry && <Button onClick={onRetry}>{__('Retry')}</Button>}
          <SelectionOptions
            options={options}
            selected={selected}
            multiple={multiple}
            label={label}
            loading={loading}
            error={error}
            hasNext={hasNext}
            dataTestId={dataTestId}
            onToggle={onToggle}
            onLoadMore={onLoadMore}
          />
          {(clearable || multiple) && (
            <Box className={styles.footer}>
              {clearable && (
                <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} onClick={onClear}>
                  {__('Clear')}
                </Button>
              )}
              {multiple && <Button onClick={onApply}>{__('Apply')}</Button>}
            </Box>
          )}
        </Box>
      </Popover>
    </Box>
  );
};
