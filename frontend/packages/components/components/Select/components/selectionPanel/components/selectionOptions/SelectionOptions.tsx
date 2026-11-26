import { _map } from '@Base/utilities/map';
import { Check, Search } from '@Components/icons';
import { Box } from '../../../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../../../Button';
import { EmptyState } from '../../../../../EmptyState';
import { Icon } from '../../../../../Icon';

import styles from '../../SelectionPanel.module.css';
import type { SelectValue } from '../../../../types/select.props';
import type { SelectionOptionsProps } from './types/selectionOptions.props';

export const SelectionOptions = <T extends SelectValue>({
  options,
  selected,
  multiple,
  label,
  loading,
  error,
  hasNext,
  dataTestId,
  onToggle,
  onLoadMore,
}: SelectionOptionsProps<T>) => (
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
        {multiple && (
          <Box
            className={styles.check}
            data-checked={selected.includes(option.value)}
            aria-hidden="true"
          >
            {selected.includes(option.value) && <Icon icon={Check} size={12} />}
          </Box>
        )}
        <Box className={styles.optionLabel}>
          {option.label}
          {option.description && (
            <Box className={styles.description}>{option.description}</Box>
          )}
        </Box>
        {!multiple && selected.includes(option.value) && <Icon icon={Check} size={14} />}
      </Button>
    ))}
    {!options.length && !loading && !error && (
      <EmptyState
        compact
        icon={Search}
        title={__('No options found')}
        description={__('Try a different search or selection.')}
      />
    )}
    {loading && <Box role="status">{__('Loading…')}</Box>}
    {hasNext && (
      <Button disabled={loading} onClick={onLoadMore}>
        {__('Load more')}
      </Button>
    )}
  </Box>
);
