import { act, renderHook } from '@testing-library/react';

import { FilterFieldType } from '../../../constants/filter.fieldType';
import { FilterOperator } from '../../../constants/filter.operator';
import { useFilters } from '../useFilters';
import type { FilterField } from '../../../types/filters.props';

const fields: readonly FilterField[] = ['email', 'status', 'createdAt'].map((id) => ({
  id,
  label: id,
  type: FilterFieldType.TEXT,
  operators: [FilterOperator.EXISTS],
}));
const first = { id: '1', field: 'email', operator: FilterOperator.EXISTS, values: [] };
const second = { id: '2', field: 'status', operator: FilterOperator.EXISTS, values: [] };

describe('useFilters', () => {
  it('excludes fields used by other rows and restores them after removal or replacement', () => {
    const { result } = renderHook(() =>
      useFilters({ fields, value: [first, second], onChange: jest.fn() }),
    );
    expect(result.current.getAvailableFields('1').map((field) => field.id)).toEqual(['email', 'createdAt']);
    expect(result.current.getAvailableFields('2').map((field) => field.id)).toEqual(['status', 'createdAt']);
    act(() => result.current.updateInline({ ...second, field: 'email' }));
    expect(result.current.draft[1]).toEqual(second);
    act(() => result.current.updateInline({ ...second, field: 'createdAt' }));
    expect(result.current.getAvailableFields('1').map((field) => field.id)).toEqual(['email', 'status']);
    act(() => result.current.remove('1'));
    expect(result.current.getAvailableFields('2').map((field) => field.id)).toEqual(['email', 'status', 'createdAt']);
  });

  it('rejects duplicate field selections in the mobile editor and blocks duplicate apply', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() => useFilters({ fields, value: [first, second], onChange }));
    act(() => result.current.edit(second));
    act(() => result.current.selectField('email'));
    expect(result.current.editing?.field).toBe('status');
    act(() => result.current.selectField('createdAt'));
    act(() => result.current.save());
    act(() => result.current.apply());
    expect(onChange).toHaveBeenCalledWith([first, { ...second, field: 'createdAt' }]);
    const duplicate = renderHook(() => useFilters({ fields, value: [first, { ...second, field: 'email' }], onChange }));
    onChange.mockClear();
    act(() => duplicate.result.current.apply());
    expect(onChange).not.toHaveBeenCalled();
  });

  it('prevents adding another filter after all allowed rows are used', () => {
    const { result } = renderHook(() =>
      useFilters({ fields, value: [first, second], maxFilters: 2, onChange: jest.fn() }),
    );
    expect(result.current.canAdd).toBe(false);
    act(() => result.current.addInline());
    act(() => result.current.add());
    expect(result.current.draft).toHaveLength(2);
    expect(result.current.editing).toBeNull();
  });
});
