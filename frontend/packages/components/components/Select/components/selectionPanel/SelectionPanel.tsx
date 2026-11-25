import { _map } from '@Base/utilities/map';
import { _filter } from '@Base/utilities/filter';
import { ChevronDown } from '@Components/icons';

import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../../../ButtonWithIcon';
import { Input } from '../../../Input';
import { Sheet } from '../../../Sheet';

import styles from './SelectionPanel.module.css';
import type { SelectValue } from '../../types/select.props';
import type { SelectionPanelProps } from './types/selectionPanel.props';

export const SelectionPanel = <T extends SelectValue>({
  options,
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
  placeholder,
  disabled,
  searchable = true,
  clearable = true,
  loading,
  error,
  dataTestId,
}: SelectionPanelProps<T>) => {
  const chosen = _filter(options, (option) => selected.includes(option.value));
  const getLabel = () => {
    if (!selected.length) return placeholder ?? __('Select an option');
    if (multiple) return __('{{count}} selected', { count: selected.length });
    return chosen[0]?.label ?? String(selected[0]);
  };

  return (
    <Box className={styles.field}>
      <ButtonWithIcon
        icon={ChevronDown}
        iconPosition={BUTTON_ICON_POSITION.RIGHT}
        variant={BUTTON_VARIANT.OUTLINE}
        fullWidth
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
      <Sheet open={open} onOpenChange={onOpenChange} title={label} dataTestId={dataTestId && `${dataTestId}-options`}>
        <Box className={styles.content}>
          {searchable && (
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
                disabled={option.disabled}
                aria-pressed={selected.includes(option.value)}
                onClick={() => onToggle(option.value)}
                className={styles.option}
                dataTestId={dataTestId && `${dataTestId}-option-${option.value}`}
              >
                <Box className={styles.optionLabel}>
                  {option.label}
                  {option.description && <Box className={styles.description}>{option.description}</Box>}
                </Box>
                {selected.includes(option.value) && <Box aria-hidden="true">✓</Box>}
              </Button>
            ))}
            {!options.length && !loading && !error && <Box>{__('No options found')}</Box>}
            {loading && <Box role="status">{__('Loading…')}</Box>}
            {hasNext && <Button disabled={loading} onClick={onLoadMore}>{__('Load more')}</Button>}
          </Box>
          <Box className={styles.footer}>
            {clearable && <Button variant={BUTTON_VARIANT.OUTLINE} onClick={onClear}>{__('Clear')}</Button>}
            {multiple && <Button onClick={onApply}>{__('Apply')}</Button>}
          </Box>
        </Box>
      </Sheet>
    </Box>
  );
};
