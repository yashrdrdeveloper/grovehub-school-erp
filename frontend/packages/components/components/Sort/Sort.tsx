import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { Plus, Trash2, ArrowUpDown } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { EmptyState } from '../EmptyState';
import { Select } from '../Select';
import { Popover } from '../Popover';

import { SortDirection } from './constants/sort.direction';
import styles from './Sort.module.css';
import type { SortProps } from './types/sort.props';
import type { useSort } from './hooks/useSort';

export const SortView = ({ fields, value, maxFields, disabled, iconOnly, dataTestId, controls }: SortProps & { readonly controls: ReturnType<typeof useSort> }) => {
  const getLabel = () => {
    if (value.length) return __('Sort ({{count}})', { count: value.length });
    return __('Sort');
  };
  return (
    <Popover open={controls.open} onOpenChange={controls.onOpenChange} title={__('Manage sort')} width={600} compact desktopScrollLayout trigger={
      <ButtonWithIcon icon={ArrowUpDown} iconSize={16} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={getLabel()} title={getLabel()} aria-expanded={controls.open} disabled={disabled} dataTestId={dataTestId} onClick={() => controls.onOpenChange(true)}>
        {!iconOnly && getLabel()}
        {iconOnly && !!value.length && value.length}
      </ButtonWithIcon>
    }>
        <Box className={styles.panel}>
          <Box className={styles.list} data-empty={!controls.draft.length}>
          {!controls.draft.length && <EmptyState compact icon={ArrowUpDown} title={__('No sorting added')} description={__('Add a sort rule to change the order of results.')} />}
          {_map(controls.draft, (rule, index) => (
            <Box className={styles.rule} key={rule.field}>
              <Select hideLabel label={__('Field')} value={rule.field} clearable={false} options={_map(_filter(fields, (field) => field.id === rule.field || !controls.draft.some((item) => item.field === field.id)), (field) => ({ value: field.id, label: field.label }))} onChange={(field) => {
                if (field) controls.update(index, { field });
              }} />
              <Select<SortDirection> hideLabel label={__('Direction')} value={rule.direction} clearable={false} searchable={false} options={[{ value: SortDirection.ASC, label: __('Ascending') }, { value: SortDirection.DESC, label: __('Descending') }]} onChange={(direction) => {
                if (direction) controls.update(index, { direction });
              }} />
              <Box className={styles.actions}>
                <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Increase priority')} disabled={index === 0} onClick={() => controls.move(index, -1)}>↑</Button>
                <Button variant={BUTTON_VARIANT.GHOST} aria-label={__('Decrease priority')} disabled={index === controls.draft.length - 1} onClick={() => controls.move(index, 1)}>↓</Button>
                <ButtonWithIcon icon={Trash2} iconSize={14} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Remove sort')} onClick={() => controls.remove(index)} />
              </Box>
            </Box>
          ))}
          </Box>
          <Box className={styles.footer}>
            <ButtonWithIcon icon={Plus} iconSize={14} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} className={styles.add} disabled={!controls.available.length || (maxFields !== undefined && controls.draft.length >= maxFields)} onClick={controls.add}>{__('Add sort')}</ButtonWithIcon>
            <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} onClick={controls.reset}>{__('Reset')}</Button>
            <Button size={BUTTON_SIZE.SMALL} onClick={controls.apply}>{__('Apply sort')}</Button>
          </Box>
        </Box>
    </Popover>
  );
};
