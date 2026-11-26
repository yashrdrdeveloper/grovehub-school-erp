import { _map } from '@Base/utilities/map';
import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { Input } from '../../../Input';
import { Select } from '../../../Select';
import { Tooltip } from '../../../Tooltip';

import { FilterValue } from '../filterValue';
import { FILTER_OPERATOR_LABELS } from '../../constants/filters.copy';
import { FilterOperator } from '../../constants/filter.operator';
import { isFilterValid } from '../../helpers/filters.isValid';
import styles from '../../Filters.module.css';
import type { MobileFilterEditorProps } from './types/mobileFilterEditor.props';

export const MobileFilterEditor = ({ controls, editing }: MobileFilterEditorProps) => {
  const { selectedField } = controls;
  return (
    <Box className={styles.panel}>
      <Box className={styles.list}>
        <Select
          label={__('Field')}
          placeholder={__('Select field')}
          value={editing.field || null}
          options={_map(controls.getAvailableFields(editing.id), (field) => ({
            value: field.id,
            label: field.label,
          }))}
          onChange={controls.selectField}
          clearable={false}
        />
        <Select<FilterOperator>
          label={__('Operator')}
          placeholder={__('Select operator')}
          value={(selectedField && editing.operator) || null}
          options={_map(selectedField?.operators ?? [], (operator) => ({
            value: operator,
            label: FILTER_OPERATOR_LABELS[operator],
          }))}
          disabled={!selectedField}
          disabledReason={__('Select a field first.')}
          clearable={false}
          onChange={(operator) => {
            if (operator)
              controls.setEditing({
                ...editing,
                operator,
                values: [],
                dateRangePresetId: undefined,
              });
          }}
        />
        {!selectedField && (
          <Tooltip label={__('Select a field first.')} focusable fullWidth>
            <Input label={__('Value')} placeholder={__('Enter value')} value="" disabled />
          </Tooltip>
        )}
        {selectedField &&
          editing.operator !== FilterOperator.EXISTS &&
          editing.operator !== FilterOperator.NOT_EXISTS && (
            <FilterValue
              field={selectedField}
              condition={editing}
              onChange={(values, dateRangePresetId) =>
                controls.setEditing({ ...editing, values, dateRangePresetId })
              }
            />
          )}
        {selectedField &&
          (editing.operator === FilterOperator.EXISTS ||
            editing.operator === FilterOperator.NOT_EXISTS) && (
            <Tooltip label={__('This operator does not require a value.')} focusable fullWidth>
              <Input label={__('Value')} placeholder={__('Not required')} value="" disabled />
            </Tooltip>
          )}
      </Box>
      <Box className={styles.footer}>
        <Button
          size={BUTTON_SIZE.SMALL}
          variant={BUTTON_VARIANT.OUTLINE}
          onClick={controls.back}
        >
          {__('Back')}
        </Button>
        <Tooltip
          label={
            (!isFilterValid(editing, selectedField) &&
              __('Select a field, operator and valid value before saving.')) ||
            undefined
          }
          focusable={Boolean(
            (!isFilterValid(editing, selectedField) &&
              __('Select a field, operator and valid value before saving.')) ||
              undefined,
          )}
        >
          <Button
            size={BUTTON_SIZE.SMALL}
            disabled={!isFilterValid(editing, selectedField)}
            onClick={controls.save}
          >
            {__('Save filter')}
          </Button>
        </Tooltip>
      </Box>
    </Box>
  );
};
