import { act, renderHook } from '@testing-library/react';

import { SortDirection } from '../../../constants/sort.direction';
import { useSort } from '../useSort';

const fields = ['email', 'status', 'createdAt'].map((id) => ({ id, label: id }));
const first = { field: 'email', direction: SortDirection.ASC };
const second = { field: 'status', direction: SortDirection.DESC };

describe('useSort', () => {
  it('excludes other rows and preserves the current field', () => {
    const { result } = renderHook(() =>
      useSort({ fields, value: [first, second], onChange: jest.fn() }),
    );
    expect(result.current.getAvailableFields(0).map((field) => field.id)).toEqual(['email', 'createdAt']);
    expect(result.current.getAvailableFields(1).map((field) => field.id)).toEqual(['status', 'createdAt']);
    act(() => result.current.update(1, { field: 'email' }));
    expect(result.current.draft[1]).toEqual(second);
    act(() => result.current.update(1, { field: 'createdAt' }));
    expect(result.current.getAvailableFields(0).map((field) => field.id)).toEqual(['email', 'status']);
    act(() => result.current.remove(0));
    expect(result.current.getAvailableFields(0).map((field) => field.id)).toEqual(['email', 'status', 'createdAt']);
  });

  it('adds only unused fields and stops when none remain', () => {
    const { result } = renderHook(() =>
      useSort({ fields, value: [first, second], onChange: jest.fn() }),
    );
    act(() => result.current.add());
    expect(result.current.draft.map((item) => item.field)).toEqual(['email', 'status', 'createdAt']);
    expect(result.current.canAdd).toBe(false);
    act(() => result.current.add());
    expect(result.current.draft).toHaveLength(3);
  });

  it('preserves sort directions and applies the dragged priority order', () => {
    const onChange = jest.fn();
    const third = { field: 'createdAt', direction: SortDirection.ASC };
    const { result } = renderHook(() =>
      useSort({ fields, value: [first, second, third], onChange }),
    );
    act(() => result.current.reorder([third, first, second]));
    expect(result.current.draft).toEqual([third, first, second]);
    expect(onChange).not.toHaveBeenCalled();
    act(() => result.current.apply());
    expect(onChange).toHaveBeenCalledWith([third, first, second]);
  });

  it('does not apply duplicate fields', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() =>
      useSort({ fields, value: [first, { ...second, field: 'email' }], onChange }),
    );
    expect(result.current.isDraftValid).toBe(false);
    act(() => result.current.apply());
    expect(onChange).not.toHaveBeenCalled();
  });
});
