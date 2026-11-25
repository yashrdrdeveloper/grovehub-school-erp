import { _map } from '@Base/utilities/map';
import { _find } from '@Base/utilities/find';
import { Filter, Plus, Trash2 } from '@Components/icons';

import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';
import { EmptyState } from '../../../EmptyState';
import { Input } from '../../../Input';
import { Select } from '../../../Select';

import { FilterValue } from '../../components/filterValue';
import { FILTER_OPERATOR_LABELS } from '../../constants/filters.copy';
import { isFilterValid } from '../../helpers/filters.isValid';
import { FilterOperator } from '../../constants/filter.operator';
import type { DesktopFiltersProps } from './types/desktopFilters.props';
import styles from './DesktopFilters.module.css';

export const DesktopFilters = ({ fields, maxFilters, controls }: DesktopFiltersProps) => {
  const invalid = controls.draft.some((condition) => !isFilterValid(condition, _find(fields, (field) => field.id === condition.field)));
  return (
    <Box className={styles.panel}>
      <Box className={styles.list}>
        {!controls.draft.length && <EmptyState compact icon={Filter} title={__('No filters added')} description={__('Add a filter to narrow down your results.')} />}
        {_map(controls.draft, (condition) => {
          const field = _find(fields, (entry) => entry.id === condition.field);
          const needsValue = condition.operator !== FilterOperator.EXISTS && condition.operator !== FilterOperator.NOT_EXISTS;
          return (
            <Box className={styles.row} key={condition.id}>
              <Select label={__('Field')} hideLabel placeholder={__('Select field')} value={condition.field || null} options={_map(fields, (entry) => ({ value: entry.id, label: entry.label }))} clearable={false} onChange={(id) => {
                const next = _find(fields, (entry) => entry.id === id);
                if (next && next.operators[0]) controls.updateInline({ ...condition, field: next.id, operator: next.operators[0], values: [] });
              }} />
              <Select<FilterOperator> label={__('Operator')} hideLabel placeholder={__('Operator')} value={(field && condition.operator) || null} options={_map(field?.operators ?? [], (operator) => ({ value: operator, label: FILTER_OPERATOR_LABELS[operator] }))} disabled={!field} searchable={false} clearable={false} onChange={(operator) => {
                if (operator) controls.updateInline({ ...condition, operator, values: [] });
              }} />
              <Box className={styles.value}>
                {field && needsValue && <FilterValue hideLabel field={field} condition={condition} onChange={(values) => controls.updateInline({ ...condition, values })} />}
                {!field && <Input label={__('Value')} placeholder={__('Enter value')} value="" disabled />}
                {field && !needsValue && <Input label={__('Value')} placeholder={__('Not required')} value="" disabled />}
              </Box>
              <ButtonWithIcon icon={Trash2} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Delete filter')} onClick={() => controls.remove(condition.id)} />
            </Box>
          );
        })}
      </Box>
      <Box className={styles.footer}>
        <ButtonWithIcon icon={Plus} iconSize={14} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} disabled={maxFilters !== undefined && controls.draft.length >= maxFilters} onClick={controls.addInline}>{__('Add filter')}</ButtonWithIcon>
        <Box className={styles.actions}>
          <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} onClick={controls.reset}>{__('Reset')}</Button>
          <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.PRIMARY} disabled={invalid} onClick={controls.apply}>{__('Save')}</Button>
        </Box>
      </Box>
    </Box>
  );
};
