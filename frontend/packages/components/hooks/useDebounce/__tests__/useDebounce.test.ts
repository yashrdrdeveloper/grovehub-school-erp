import { act, renderHook } from '@testing-library/react';

import { useDebounce } from '../useDebounce';

describe('useDebounce', () => {
  beforeEach(() => jest.useFakeTimers());
  afterEach(() => jest.useRealTimers());

  it('dispatches only the last value and cancels on unmount', () => {
    const callback = jest.fn();
    const { result, unmount } = renderHook(() => useDebounce(callback, 350));

    act(() => {
      result.current.execute('a');
      result.current.execute('ab');
      jest.advanceTimersByTime(349);
    });
    expect(callback).not.toHaveBeenCalled();
    act(() => jest.advanceTimersByTime(1));
    expect(callback).toHaveBeenCalledWith('ab');
    act(() => result.current.execute('abc'));
    unmount();
    act(() => jest.runOnlyPendingTimers());
    expect(callback).toHaveBeenCalledTimes(1);
  });

  it('flushes once and cancels pending work', () => {
    const callback = jest.fn();
    const { result } = renderHook(() => useDebounce(callback, 350));
    act(() => {
      result.current.execute('a');
      result.current.flush();
      result.current.flush();
      result.current.execute('b');
      result.current.cancel();
      jest.runOnlyPendingTimers();
    });
    expect(callback).toHaveBeenCalledTimes(1);
    expect(callback).toHaveBeenCalledWith('a');
  });
});
