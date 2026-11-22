import { act, renderHook } from '@testing-library/react';
import { useBoolean } from '../useBoolean';

describe('useBoolean', () => {
  it('supports predictable state transitions with stable callbacks', () => {
    const { result, rerender } = renderHook(() => useBoolean());
    expect(result.current.value).toBe(false);
    const setTrue = result.current.setTrue;
    act(() => result.current.setTrue());
    expect(result.current.value).toBe(true);
    act(() => result.current.toggle());
    expect(result.current.value).toBe(false);
    act(() => result.current.setValue(true));
    expect(result.current.value).toBe(true);
    act(() => result.current.setFalse());
    rerender();
    expect(result.current.value).toBe(false);
    expect(result.current.setTrue).toBe(setTrue);
  });

  it('honors an initial true value', () => {
    const { result } = renderHook(() => useBoolean(true));
    expect(result.current.value).toBe(true);
  });
});
