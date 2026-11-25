import { _map } from '@Base/utilities/map';
import { Plus, Trash2, ArrowUpDown } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { EmptyState } from '../EmptyState';
import { Select } from '../Select';
import { Tooltip } from '../Tooltip';
import { Popover } from '../Popover';

import { SortDirection } from './constants/sort.direction';
import styles from './Sort.module.css';
import type { SortProps } from './types/sort.props';
import type { useSort } from './hooks/useSort';

export const SortView = ({ value, maxFields, disabled, iconOnly, dataTestId, controls }: SortProps & { readonly controls: ReturnType<typeof useSort> }) => {
  const getLabel = () => {
    if (value.length) return __('Sort ({{count}})', { count: value.length });
    return __('Sort');
  };
  return (
    <Popover open={controls.open} onOpenChange={controls.onOpenChange} title={__('Manage sort')} width={600} compact desktopScrollLayout trigger={
      <Tooltip label={getLabel()}><ButtonWithIcon icon={ArrowUpDown} iconSize={16} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={getLabel()} aria-expanded={controls.open} disabled={disabled} dataTestId={dataTestId} onClick={() => controls.onOpenChange(true)}>
        {!iconOnly && getLabel()}
        {iconOnly && !!value.length && value.length}
      </ButtonWithIcon></Tooltip>
    }>
        <Box className={styles.panel}>
          <Box className={styles.list} data-empty={!controls.draft.length}>
          {!controls.draft.length && <EmptyState compact icon={ArrowUpDown} title={__('No sorting added')} description={__('Add a sort rule to change the order of results.')} />}
          {_map(controls.draft, (rule, index) => (
            <Box className={styles.rule} key={rule.field}>
              <Select hideLabel label={__('Field')} value={rule.field} clearable={false} options={_map(controls.getAvailableFields(index), (field) => ({ value: field.id, label: field.label }))} onChange={(field) => {
                if (field) controls.update(index, { field });
              }} />
              <Select<SortDirection> hideLabel label={__('Direction')} value={rule.direction} clearable={false} searchable={false} options={[{ value: SortDirection.ASC, label: __('Ascending') }, { value: SortDirection.DESC, label: __('Descending') }]} onChange={(direction) => {
                if (direction) controls.update(index, { direction });
              }} />
              <Box className={styles.actions}>
                {maxFields !== 1 && controls.draft.length > 1 && (
                  <>
                    <Tooltip label={index === 0 && __('This sort already has the highest priority.') || undefined} focusable={Boolean(index === 0 && __('This sort already has the highest priority.') || undefined)}><Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Increase priority')} disabled={index === 0} onClick={() => controls.move(index, -1)}>↑</Button></Tooltip>
                    <Tooltip label={index === controls.draft.length - 1 && __('This sort already has the lowest priority.') || undefined} focusable={Boolean(index === controls.draft.length - 1 && __('This sort already has the lowest priority.') || undefined)}><Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Decrease priority')} disabled={index === controls.draft.length - 1} onClick={() => controls.move(index, 1)}>↓</Button></Tooltip>
                  </>
                )}
                <ButtonWithIcon icon={Trash2} iconSize={14} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Remove sort')} onClick={() => controls.remove(index)} />
              </Box>
            </Box>
          ))}
          </Box>
          <Box className={styles.footer}>
            <Tooltip label={controls.addReason} focusable={Boolean(controls.addReason)}><ButtonWithIcon icon={Plus} iconSize={14} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} className={styles.add} disabled={!controls.canAdd} onClick={controls.add}>{__('Add sort')}</ButtonWithIcon></Tooltip>
            <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} onClick={controls.reset}>{__('Reset')}</Button>
            <Tooltip label={!controls.isDraftValid && __('Choose a unique field for each sort rule.') || undefined} focusable={Boolean(!controls.isDraftValid && __('Choose a unique field for each sort rule.') || undefined)}><Button size={BUTTON_SIZE.SMALL} disabled={!controls.isDraftValid} onClick={controls.apply}>{__('Apply sort')}</Button></Tooltip>
          </Box>
        </Box>
    </Popover>
  );
};
