import { _map } from '@Base/utilities/map';
import { _find } from '@Base/utilities/find';
import { Filter, Plus, Trash2 } from '@Components/icons';

import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';
import { EmptyState } from '../../../EmptyState';
import { Input } from '../../../Input';
import { Select } from '../../../Select';
import { Tooltip } from '../../../Tooltip';

import { FilterValue } from '../../components/filterValue';
import { FILTER_OPERATOR_LABELS } from '../../constants/filters.copy';
import { FilterOperator } from '../../constants/filter.operator';
import type { DesktopFiltersProps } from './types/desktopFilters.props';
import styles from './DesktopFilters.module.css';

export const DesktopFilters = ({ fields, controls }: DesktopFiltersProps) => {
  return (
    <Box className={styles.panel}>
      <Box className={styles.list} data-empty={!controls.draft.length}>
        {!controls.draft.length && <EmptyState compact icon={Filter} title={__('No filters added')} description={__('Add a filter to narrow down your results.')} />}
        {_map(controls.draft, (condition) => {
          const field = _find(fields, (entry) => entry.id === condition.field);
          const needsValue = condition.operator !== FilterOperator.EXISTS && condition.operator !== FilterOperator.NOT_EXISTS;
          return (
            <Box className={styles.row} key={condition.id}>
              <Select label={__('Field')} hideLabel placeholder={__('Select field')} value={condition.field || null} options={_map(controls.getAvailableFields(condition.id), (entry) => ({ value: entry.id, label: entry.label }))} clearable={false} onChange={(id) => {
                const next = _find(controls.getAvailableFields(condition.id), (entry) => entry.id === id);
                if (next && next.operators[0]) controls.updateInline({ ...condition, field: next.id, operator: next.operators[0], values: [] });
              }} />
              <Select<FilterOperator> label={__('Operator')} hideLabel placeholder={__('Select operator')} value={(field && condition.operator) || null} options={_map(field?.operators ?? [], (operator) => ({ value: operator, label: FILTER_OPERATOR_LABELS[operator] }))} disabled={!field} disabledReason={__('Select a field first.')} searchable={false} clearable={false} onChange={(operator) => {
                if (operator) controls.updateInline({ ...condition, operator, values: [] });
              }} />
              <Box className={styles.value}>
                {field && needsValue && <FilterValue hideLabel field={field} condition={condition} onChange={(values) => controls.updateInline({ ...condition, values })} />}
                {!field && <Tooltip label={__('Select a field first.')} focusable fullWidth><Input label={__('Value')} placeholder={__('Enter value')} value="" disabled /></Tooltip>}
                {field && !needsValue && <Tooltip label={__('This operator does not require a value.')} focusable fullWidth><Input label={__('Value')} placeholder={__('Not required')} value="" disabled /></Tooltip>}
              </Box>
              <ButtonWithIcon icon={Trash2} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Delete filter')} onClick={() => controls.remove(condition.id)} />
            </Box>
          );
        })}
      </Box>
      <Box className={styles.footer}>
        <Tooltip label={controls.addReason} focusable={Boolean(controls.addReason)}><ButtonWithIcon icon={Plus} iconSize={14} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} disabled={!controls.canAdd} onClick={controls.addInline}>{__('Add filter')}</ButtonWithIcon></Tooltip>
        <Box className={styles.actions}>
          <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} onClick={controls.reset}>{__('Reset')}</Button>
          <Tooltip label={!controls.isDraftValid && __('Complete all filter fields and values before applying.') || undefined} focusable={Boolean(!controls.isDraftValid && __('Complete all filter fields and values before applying.') || undefined)}><Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.PRIMARY} disabled={!controls.isDraftValid} onClick={controls.apply}>{__('Apply filters')}</Button></Tooltip>
        </Box>
      </Box>
    </Box>
  );
};
