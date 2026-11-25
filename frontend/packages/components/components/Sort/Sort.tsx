import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { Plus, Trash2, Settings2 } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Select } from '../Select';
import { Sheet } from '../Sheet';

import { SortDirection } from './constants/sort.direction';
import styles from './Sort.module.css';
import type { SortProps } from './types/sort.props';
import type { useSort } from './hooks/useSort';

export const SortView = ({ fields, value, maxFields, disabled, dataTestId, controls }: SortProps & { readonly controls: ReturnType<typeof useSort> }) => {
  const getLabel = () => {
    if (value.length) return __('Sort ({{count}})', { count: value.length });
    return __('Sort');
  };
  return (
    <>
      <ButtonWithIcon icon={Settings2} variant={BUTTON_VARIANT.OUTLINE} disabled={disabled} dataTestId={dataTestId} onClick={() => controls.onOpenChange(true)}>{getLabel()}</ButtonWithIcon>
      <Sheet open={controls.open} onOpenChange={controls.onOpenChange} title={__('Manage sort')}>
        <Box className={styles.panel}>
          {_map(controls.draft, (rule, index) => (
            <Box className={styles.rule} key={rule.field}>
              <Select label={__('Field')} value={rule.field} clearable={false} options={_map(_filter(fields, (field) => field.id === rule.field || !controls.draft.some((item) => item.field === field.id)), (field) => ({ value: field.id, label: field.label }))} onChange={(field) => {
                if (field) controls.update(index, { field });
              }} />
              <Select<SortDirection> label={__('Direction')} value={rule.direction} clearable={false} searchable={false} options={[{ value: SortDirection.ASC, label: __('Ascending') }, { value: SortDirection.DESC, label: __('Descending') }]} onChange={(direction) => {
                if (direction) controls.update(index, { direction });
              }} />
              <Box className={styles.actions}>
                <Button variant={BUTTON_VARIANT.GHOST} aria-label={__('Increase priority')} disabled={index === 0} onClick={() => controls.move(index, -1)}>↑</Button>
                <Button variant={BUTTON_VARIANT.GHOST} aria-label={__('Decrease priority')} disabled={index === controls.draft.length - 1} onClick={() => controls.move(index, 1)}>↓</Button>
                <ButtonWithIcon icon={Trash2} variant={BUTTON_VARIANT.GHOST} aria-label={__('Remove sort')} onClick={() => controls.remove(index)} />
              </Box>
            </Box>
          ))}
          <ButtonWithIcon icon={Plus} variant={BUTTON_VARIANT.OUTLINE} disabled={!controls.available.length || (maxFields !== undefined && controls.draft.length >= maxFields)} onClick={controls.add}>{__('Add sort')}</ButtonWithIcon>
          <Box className={styles.actions}>
            <Button variant={BUTTON_VARIANT.OUTLINE} onClick={controls.reset}>{__('Reset')}</Button>
            <Button onClick={controls.apply}>{__('Apply sort')}</Button>
          </Box>
        </Box>
      </Sheet>
    </>
  );
};
