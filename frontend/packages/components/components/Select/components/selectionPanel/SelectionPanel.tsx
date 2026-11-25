import { _map } from '@Base/utilities/map';
import { _filter } from '@Base/utilities/filter';
import { ChevronDown, Check, Search } from '@Components/icons';

import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../../../ButtonWithIcon';
import { EmptyState } from '../../../EmptyState';
import { Input } from '../../../Input';
import { InputControl, INPUT_CONTROL_SIZE } from '../../../InputControl';
import { Icon } from '../../../Icon';
import { Popover } from '../../../Popover';

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
    return undefined;
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
      <Popover open={open} onOpenChange={onOpenChange} title={label} compact fullWidth desktopDropdown desktopDropdownWidth={getDropdownWidth()} dataTestId={dataTestId && `${dataTestId}-options`} trigger={
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
        onClick={() => onOpenChange(true)}
      >
        {getLabel()}
      </ButtonWithIcon>
      }>
        <Box className={styles.content} data-dropdown={dropdown} data-multiple={multiple}>
          {searchable && dropdown && (
            <Box className={styles.search}>
              <Box className={styles.searchIcon} aria-hidden="true"><Icon icon={Search} size={14} /></Box>
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
          <Box className={styles.options} role="group" aria-label={label}>
            {_map(options, (option) => (
              <Button
                key={option.value}
                variant={BUTTON_VARIANT.GHOST}
                size={BUTTON_SIZE.SMALL}
                disabled={option.disabled}
                aria-pressed={selected.includes(option.value)}
                onClick={() => onToggle(option.value)}
                className={styles.option}
                dataTestId={dataTestId && `${dataTestId}-option-${option.value}`}
              >
                {multiple && <Box className={styles.check} data-checked={selected.includes(option.value)} aria-hidden="true">{selected.includes(option.value) && <Icon icon={Check} size={12} />}</Box>}
                <Box className={styles.optionLabel}>
                  {option.label}
                  {option.description && <Box className={styles.description}>{option.description}</Box>}
                </Box>
                {!multiple && selected.includes(option.value) && <Icon icon={Check} size={14} />}
              </Button>
            ))}
            {!options.length && !loading && !error && <EmptyState compact icon={Search} title={__('No options found')} description={__('Try a different search or selection.')} />}
            {loading && <Box role="status">{__('Loading…')}</Box>}
            {hasNext && <Button disabled={loading} onClick={onLoadMore}>{__('Load more')}</Button>}
          </Box>
          {(clearable || multiple) && (
            <Box className={styles.footer}>
              {clearable && <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} onClick={onClear}>{__('Clear')}</Button>}
              {multiple && <Button onClick={onApply}>{__('Apply')}</Button>}
            </Box>
          )}
        </Box>
      </Popover>
    </Box>
  );
};
