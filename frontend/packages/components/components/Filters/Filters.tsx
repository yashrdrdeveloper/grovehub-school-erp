import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { _map } from '@Base/utilities/map';
import { _find } from '@Base/utilities/find';
import { Pencil, Plus, Trash2, Filter } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Select } from '../Select';
import { Popover } from '../Popover';

import { DesktopFilters } from './templates/desktopFilters';
import { FilterValue } from './components/filterValue';
import { FILTER_OPERATOR_LABELS } from './constants/filters.copy';
import { isFilterValid } from './helpers/filters.isValid';
import styles from './Filters.module.css';
import type { FilterOperator } from './constants/filter.operator';
import type { FiltersProps } from './types/filters.props';
import type { useFilters } from './hooks/useFilters';

export const FiltersView = ({ fields, value, disabled, maxFilters, iconOnly, dataTestId, controls }: FiltersProps & { readonly controls: ReturnType<typeof useFilters> }) => {
  const mobile = useMobileViewport();
  const { editing, selectedField } = controls;
  const getTitle = () => {
    if (!mobile) return __('Manage filters');
    if (editing) return __('Filter condition');
    return __('Manage filters');
  };
  const getButton = () => {
    if (value.length) return __('Filters ({{count}})', { count: value.length });
    return __('Filters');
  };
  return (
    <Popover open={controls.open} onOpenChange={controls.onOpenChange} title={getTitle()} width={520} trigger={
      <ButtonWithIcon icon={Filter} size={BUTTON_SIZE.SMALL} iconSize={16} variant={BUTTON_VARIANT.OUTLINE} aria-label={getButton()} aria-expanded={controls.open} title={getButton()} disabled={disabled} dataTestId={dataTestId} onClick={() => controls.onOpenChange(true)}>
        {(!iconOnly || !!value.length) && getButton()}
      </ButtonWithIcon>
    }>
        {!mobile && <DesktopFilters fields={fields} maxFilters={maxFilters} controls={controls} />}
        {mobile && !editing && (
          <Box className={styles.panel}>
            <Box className={styles.list}>
              {!controls.draft.length && <Box>{__('No filters added')}</Box>}
              {_map(controls.draft, (condition) => (
                <Box className={styles.condition} key={condition.id}>
                  <Box className={styles.summary}>
                    <Box>{_find(fields, (field) => field.id === condition.field)?.label ?? condition.field}</Box>
                    <Box className={styles.muted}>{FILTER_OPERATOR_LABELS[condition.operator]} {condition.values.join(', ')}</Box>
                  </Box>
                  <ButtonWithIcon icon={Pencil} variant={BUTTON_VARIANT.GHOST} aria-label={__('Edit filter')} onClick={() => controls.edit(condition)} />
                  <ButtonWithIcon icon={Trash2} variant={BUTTON_VARIANT.GHOST} aria-label={__('Delete filter')} onClick={() => controls.remove(condition.id)} />
                </Box>
              ))}
              <ButtonWithIcon icon={Plus} variant={BUTTON_VARIANT.OUTLINE} disabled={maxFilters !== undefined && controls.draft.length >= maxFilters} onClick={controls.add}>{__('Add filter')}</ButtonWithIcon>
            </Box>
            <Box className={styles.footer}>
              <Button variant={BUTTON_VARIANT.OUTLINE} onClick={controls.reset}>{__('Reset')}</Button>
              <Button onClick={controls.apply}>{__('Apply filters')}</Button>
            </Box>
          </Box>
        )}
        {mobile && editing && (
          <Box className={styles.panel}>
            <Box className={styles.list}>
              <Select label={__('Field')} value={editing.field || null} options={_map(fields, (field) => ({ value: field.id, label: field.label }))} onChange={controls.selectField} clearable={false} />
              {selectedField && (
                <>
                  <Select<FilterOperator> label={__('Operator')} value={editing.operator} options={_map(selectedField.operators, (operator) => ({ value: operator, label: FILTER_OPERATOR_LABELS[operator] }))} clearable={false} onChange={(operator) => {
                    if (operator) controls.setEditing({ ...editing, operator, values: [] });
                  }} />
                  <FilterValue field={selectedField} condition={editing} onChange={(values) => controls.setEditing({ ...editing, values })} />
                </>
              )}
            </Box>
            <Box className={styles.footer}>
              <Button variant={BUTTON_VARIANT.OUTLINE} onClick={controls.back}>{__('Back')}</Button>
              <Button disabled={!isFilterValid(editing, selectedField)} onClick={controls.save}>{__('Save filter')}</Button>
            </Box>
          </Box>
        )}
    </Popover>
  );
};
