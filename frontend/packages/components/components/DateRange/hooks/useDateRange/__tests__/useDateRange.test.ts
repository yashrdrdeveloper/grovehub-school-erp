import { act, renderHook } from '@testing-library/react';

import { useDateRange } from '../useDateRange';

describe('useDateRange', () => {
  it('keeps incomplete ranges in the draft and commits only when applied', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() => useDateRange({ value: { startDate: null, endDate: null }, onChange }));
    act(() => result.current.onOpenChange(true));
    act(() => result.current.selectDate('2026-11-10'));
    expect(result.current.canApply).toBe(false);
    act(() => result.current.apply());
    expect(onChange).not.toHaveBeenCalled();
    act(() => result.current.setHover('2026-11-15'));
    expect(result.current.preview.endDate).toBe('2026-11-15');
    expect(result.current.draft.endDate).toBeNull();
    act(() => result.current.selectDate('2026-11-15'));
    act(() => result.current.apply());
    expect(onChange).toHaveBeenCalledWith({ startDate: '2026-11-10', endDate: '2026-11-15' });
    expect(result.current.open).toBe(false);
  });
});
