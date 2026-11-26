import { AsyncMultiSelect } from '../../../AsyncMultiSelect';
import { AsyncSelect } from '../../../AsyncSelect';
import { Date as DateInput } from '../../../Date';
import { DateTime } from '../../../DateTime';
import { DateRange } from '../../../DateRange';
import { Input } from '../../../Input';
import { MultiSelect } from '../../../MultiSelect';
import { Number as NumberInput } from '../../../Number';
import { NumberRange } from '../../../NumberRange';
import { Select } from '../../../Select';

import { FilterFieldType } from '../../constants/filter.fieldType';
import { FilterOperator } from '../../constants/filter.operator';
import type { FilterValueProps } from './types/filterValue.props';

export const FilterValue = ({ field, condition, onChange, hideLabel }: FilterValueProps) => {
  const first = condition.values[0] ?? '';
  const second = condition.values[1] ?? '';
  const label = __('Value');
  const single = (value: string | number | null) => {
    if (value === null) {
      onChange([]);
      return;
    }
    onChange([String(value)]);
  };
  const numeric = (value: string) => {
    if (!value) return null;
    const result = globalThis.Number(value);
    if (globalThis.Number.isFinite(result)) return result;
    return null;
  };
  const multiple =
    condition.operator === FilterOperator.IN || condition.operator === FilterOperator.NOT_IN;
  if (
    condition.operator === FilterOperator.EXISTS ||
    condition.operator === FilterOperator.NOT_EXISTS
  )
    return null;
  if (field.type === FilterFieldType.CUSTOM)
    return field.renderValue?.(condition, onChange) ?? null;
  if (
    field.type === FilterFieldType.ASYNC_SELECT ||
    field.type === FilterFieldType.ASYNC_MULTI_SELECT
  ) {
    if (!field.source) return null;
    if (multiple || field.type === FilterFieldType.ASYNC_MULTI_SELECT) {
      return (
        <AsyncMultiSelect
          {...field.source}
          label={label}
          placeholder={__('Select value')}
          hideLabel={hideLabel}
          value={condition.values}
          onChange={onChange}
        />
      );
    }
    return (
      <AsyncSelect
        {...field.source}
        label={label}
        placeholder={__('Select value')}
        hideLabel={hideLabel}
        value={first || null}
        onChange={single}
      />
    );
  }
  if (
    field.type === FilterFieldType.SELECT ||
    field.type === FilterFieldType.MULTI_SELECT ||
    field.type === FilterFieldType.BOOLEAN
  ) {
    if (multiple || field.type === FilterFieldType.MULTI_SELECT) {
      return (
        <MultiSelect
          label={label}
          placeholder={__('Select value')}
          hideLabel={hideLabel}
          options={field.options ?? []}
          value={condition.values}
          onChange={onChange}
        />
      );
    }
    return (
      <Select
        label={label}
        placeholder={__('Select value')}
        hideLabel={hideLabel}
        options={field.options ?? []}
        value={first || null}
        onChange={single}
      />
    );
  }
  if (
    field.type === FilterFieldType.DATE_RANGE ||
    (field.type === FilterFieldType.DATE && condition.operator === FilterOperator.BETWEEN)
  ) {
    return (
      <DateRange
        showHeader={false}
        showFooter={false}
        value={{
          startDate: first || null,
          endDate: second || null,
          presetId: condition.dateRangePresetId,
        }}
        onChange={(range) => onChange([range.startDate ?? '', range.endDate ?? ''], range.presetId)}
      />
    );
  }
  if (
    field.type === FilterFieldType.NUMBER_RANGE ||
    (field.type === FilterFieldType.NUMBER && condition.operator === FilterOperator.BETWEEN)
  ) {
    return (
      <NumberRange
        value={{ min: numeric(first), max: numeric(second) }}
        onChange={(range) => onChange([String(range.min ?? ''), String(range.max ?? '')])}
      />
    );
  }
  if (field.type === FilterFieldType.NUMBER)
    return <NumberInput label={label} value={numeric(first)} onChange={single} />;
  if (field.type === FilterFieldType.DATE)
    return <DateInput label={label} value={first || null} onChange={single} />;
  if (field.type === FilterFieldType.DATE_TIME)
    return <DateTime label={label} value={first || null} onChange={single} />;
  return (
    <Input
      label={label}
      placeholder={__('Enter value')}
      value={first}
      onChange={(event) => single(event.currentTarget.value)}
    />
  );
};
