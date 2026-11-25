import { act, renderHook } from '@testing-library/react';

import { DATE_RANGE_PRESETS } from '../../../constants/dateRange.presets';
import type { DateRangeValue } from '../../../types/dateRange.props';
import { useDateRange } from '../useDateRange';

describe('useDateRange', () => {
  it('commits a completed manual range when clicking outside', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() => useDateRange({ value: { startDate: null, endDate: null }, onChange }));
    act(() => result.current.onOpenChange(true));
    act(() => result.current.selectDate('2026-11-10'));
    act(() => result.current.selectDate('2026-11-15'));
    act(() => result.current.onOutsideClick());
    expect(onChange).toHaveBeenCalledTimes(1);
    expect(onChange).toHaveBeenCalledWith({ startDate: '2026-11-10', endDate: '2026-11-15', presetId: 'custom' });
    expect(result.current.open).toBe(false);
  });

  it('commits a preset and retains its identity on outside dismissal', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() => useDateRange({ value: { startDate: null, endDate: null }, onChange }));
    act(() => result.current.onOpenChange(true));
    act(() => result.current.selectPreset(DATE_RANGE_PRESETS[0]));
    act(() => result.current.onOutsideClick());
    expect(onChange).toHaveBeenCalledWith(expect.objectContaining({ presetId: DATE_RANGE_PRESETS[0].id }));
    expect(result.current.open).toBe(false);
  });

  it('discards incomplete dates and skips unchanged values on outside dismissal', () => {
    const onChange = jest.fn();
    const value = { startDate: '2026-11-01', endDate: '2026-11-30', presetId: 'thisMonth' };
    const { result } = renderHook(() => useDateRange({ value, onChange }));
    act(() => result.current.onOpenChange(true));
    act(() => result.current.onOutsideClick());
    expect(onChange).not.toHaveBeenCalled();
    act(() => result.current.onOpenChange(true));
    act(() => result.current.selectDate('2026-11-12'));
    act(() => result.current.onOutsideClick());
    expect(onChange).not.toHaveBeenCalled();
    expect(result.current.open).toBe(false);
    act(() => result.current.onOpenChange(true));
    expect(result.current.selectedPreset).toBe('thisMonth');
    expect(result.current.draft).toEqual(value);
  });

  it('restores the applied preset after closing and remounting', () => {
    const onChange = jest.fn();
    const initial: DateRangeValue = { startDate: null, endDate: null };
    const preset = DATE_RANGE_PRESETS.find((item) => item.id === 'last30Days')!;
    const first = renderHook(() => useDateRange({ value: initial, onChange }));
    act(() => first.result.current.onOpenChange(true));
    act(() => first.result.current.selectPreset(preset));
    act(() => first.result.current.apply());
    const applied = onChange.mock.calls[0][0] as DateRangeValue;
    expect(applied.presetId).toBe('last30Days');
    first.unmount();
    const second = renderHook(() => useDateRange({ value: applied, onChange }));
    act(() => second.result.current.onOpenChange(true));
    expect(second.result.current.selectedPreset).toBe('last30Days');
    act(() => second.result.current.selectCustom());
    expect(second.result.current.selectedPreset).toBe('custom');
    act(() => second.result.current.onOpenChange(false));
    act(() => second.result.current.onOpenChange(true));
    expect(second.result.current.selectedPreset).toBe('last30Days');
  });

  it('keeps manual ranges custom even when their dates match a preset', () => {
    const onChange = jest.fn();
    const value: DateRangeValue = { startDate: '2026-11-01', endDate: '2026-11-30' };
    const { result } = renderHook(() => useDateRange({ value, onChange }));
    act(() => result.current.onOpenChange(true));
    expect(result.current.selectedPreset).toBe('custom');
    act(() => result.current.selectCustom());
    act(() => result.current.apply());
    expect(onChange).toHaveBeenCalledWith({ ...value, presetId: 'custom' });
  });

  it('clears the selected preset along with the dates', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() => useDateRange({ value: { startDate: '2026-11-01', endDate: '2026-11-30', presetId: 'thisMonth' }, onChange }));
    act(() => result.current.onOpenChange(true));
    act(() => result.current.clear());
    expect(result.current.selectedPreset).toBeNull();
    act(() => result.current.apply());
    expect(onChange).toHaveBeenCalledWith({ startDate: null, endDate: null });
  });

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
    expect(onChange).toHaveBeenCalledWith({ startDate: '2026-11-10', endDate: '2026-11-15', presetId: 'custom' });
    expect(result.current.open).toBe(false);
  });
});
