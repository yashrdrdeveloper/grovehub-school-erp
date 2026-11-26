import { _map } from '@Base/utilities/map';
import { Plus, Trash2, ArrowUpDown } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { DragList } from '../DragList';
import { EmptyState } from '../EmptyState';
import { Select } from '../Select';
import { Tooltip } from '../Tooltip';
import { Popover } from '../Popover';

import { SortDirection } from './constants/sort.direction';
import styles from './Sort.module.css';
import type { SortProps } from './types/sort.props';
import type { useSort } from './hooks/useSort';

export const SortView = ({
  value,
  maxFields,
  disabled,
  iconOnly,
  dataTestId,
  initialOpen,
  controls,
}: SortProps & { readonly controls: ReturnType<typeof useSort> }) => {
  const getLabel = () => {
    if (value.length) return __('Sort ({{count}})', { count: value.length });
    return __('Sort');
  };
  const invalidReason =
    (!controls.isDraftValid && __('Choose a unique field for each sort rule.')) || undefined;

  return (
    <Popover
      open={controls.open}
      onOpenChange={controls.onOpenChange}
      title={__('Manage sort')}
      width={600}
      compact
      desktopScrollLayout
      trigger={
        !initialOpen && (
          <Tooltip label={getLabel()}>
            <ButtonWithIcon
              icon={ArrowUpDown}
              iconSize={16}
              size={BUTTON_SIZE.SMALL}
              variant={BUTTON_VARIANT.OUTLINE}
              aria-label={getLabel()}
              aria-expanded={controls.open}
              disabled={disabled}
              dataTestId={dataTestId}
              onClick={() => controls.onOpenChange(true)}
            >
              {!iconOnly && getLabel()}
              {iconOnly && !!value.length && value.length}
            </ButtonWithIcon>
          </Tooltip>
        )
      }
    >
      <Box className={styles.panel}>
        <Box className={styles.list} data-empty={!controls.draft.length}>
          {!controls.draft.length && (
            <EmptyState
              compact
              icon={ArrowUpDown}
              title={__('No sorting added')}
              description={__('Add a sort rule to change the order of results.')}
            />
          )}
          {!!controls.draft.length && (
            <DragList
              compact
              items={controls.draft}
              getItemId={(rule) => rule.field}
              label={__('Sort priority')}
              disabled={maxFields === 1 || controls.draft.length < 2}
              disabledReason={__('Add another sort rule to change its priority.')}
              onReorder={controls.reorder}
              renderItem={(rule, index) => (
                <Box className={styles.rule}>
                  <Select
                    hideLabel
                    label={__('Field')}
                    value={rule.field}
                    clearable={false}
                    options={_map(controls.getAvailableFields(index), (field) => ({
                      value: field.id,
                      label: field.label,
                    }))}
                    onChange={(field) => {
                      if (field) controls.update(index, { field });
                    }}
                  />
                  <Select<SortDirection>
                    hideLabel
                    label={__('Direction')}
                    value={rule.direction}
                    clearable={false}
                    searchable={false}
                    options={[
                      { value: SortDirection.ASC, label: __('Ascending') },
                      { value: SortDirection.DESC, label: __('Descending') },
                    ]}
                    onChange={(direction) => {
                      if (direction) controls.update(index, { direction });
                    }}
                  />
                  <Tooltip label={__('Remove sort')}>
                    <ButtonWithIcon
                      icon={Trash2}
                      iconSize={14}
                      size={BUTTON_SIZE.SMALL}
                      variant={BUTTON_VARIANT.GHOST}
                      aria-label={__('Remove sort')}
                      onClick={() => controls.remove(index)}
                    />
                  </Tooltip>
                </Box>
              )}
            />
          )}
        </Box>
        <Box className={styles.footer}>
          <Tooltip label={controls.addReason} focusable={Boolean(controls.addReason)}>
            <ButtonWithIcon
              icon={Plus}
              iconSize={14}
              size={BUTTON_SIZE.SMALL}
              variant={BUTTON_VARIANT.GHOST}
              className={styles.add}
              disabled={!controls.canAdd}
              onClick={controls.add}
            >
              {__('Add sort')}
            </ButtonWithIcon>
          </Tooltip>
          <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} onClick={controls.reset}>
            {__('Reset')}
          </Button>
          <Tooltip label={invalidReason} focusable={Boolean(invalidReason)}>
            <Button
              size={BUTTON_SIZE.SMALL}
              disabled={!controls.isDraftValid}
              onClick={controls.apply}
            >
              {__('Apply sort')}
            </Button>
          </Tooltip>
        </Box>
      </Box>
    </Popover>
  );
};
