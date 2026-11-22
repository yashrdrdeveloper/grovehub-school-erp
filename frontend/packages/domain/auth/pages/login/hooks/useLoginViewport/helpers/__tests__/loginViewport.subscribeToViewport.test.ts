import { STYLE_MEDIA_QUERY } from '@Styles/tokens/breakpoints';
import { subscribeToViewport } from '../loginViewport.subscribeToViewport';

describe('subscribeToViewport', () => {
  const initial = window.matchMedia;

  afterAll(() => {
    Object.defineProperty(window, 'matchMedia', { configurable: true, value: initial });
  });

  it('registers and removes the same browser listener', () => {
    const add = jest.fn();
    const remove = jest.fn();
    const query = jest.fn(() => ({
      matches: false,
      addEventListener: add,
      removeEventListener: remove,
    }));
    Object.defineProperty(window, 'matchMedia', { configurable: true, value: query });
    const listener = jest.fn();
    const unsubscribe = subscribeToViewport(listener);
    expect(query).toHaveBeenCalledWith(STYLE_MEDIA_QUERY.PHONE);
    expect(add).toHaveBeenCalledWith('change', listener);
    unsubscribe();
    expect(remove).toHaveBeenCalledWith('change', listener);
  });

  it('uses a safe no-op unsubscribe when matchMedia is unavailable', () => {
    Object.defineProperty(window, 'matchMedia', { configurable: true, value: undefined });
    expect(() => subscribeToViewport(() => {})()).not.toThrow();
  });
});
