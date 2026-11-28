import { isFilterValid } from '../filters.isValid';
import { FilterFieldType } from '../../constants/filter.fieldType';
import { FilterOperator } from '../../constants/filter.operator';
import type { FilterField } from '../../types/filters.props';

const field: FilterField = {
  id: 'email',
  label: 'Email',
  type: FilterFieldType.TEXT,
  operators: [FilterOperator.EQUALS, FilterOperator.EXISTS],
};

describe('isFilterValid', () => {
  it('rejects empty or unsupported values', () => {
    expect(
      isFilterValid(
        { id: '1', field: 'email', operator: FilterOperator.EQUALS, values: [''] },
        field,
      ),
    ).toBe(false);
    expect(
      isFilterValid({ id: '1', field: 'email', operator: FilterOperator.IN, values: ['a'] }, field),
    ).toBe(false);
  });

  it('accepts existence without a value', () => {
    expect(
      isFilterValid(
        { id: '1', field: 'email', operator: FilterOperator.EXISTS, values: [] },
        field,
      ),
    ).toBe(true);
  });
});
