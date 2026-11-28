import { STYLE_MEDIA_QUERY } from '@Styles/tokens/breakpoints';
import { getMobileSnapshot } from '../loginViewport.getMobileSnapshot';

describe('getMobileSnapshot', () => {
  const initial = window.matchMedia;

  afterAll(() => {
    Object.defineProperty(window, 'matchMedia', { configurable: true, value: initial });
  });

  it('reads the canonical mobile query from the browser', () => {
    const query = jest.fn(() => ({ matches: true }));
    Object.defineProperty(window, 'matchMedia', { configurable: true, value: query });
    expect(getMobileSnapshot()).toBe(true);
    expect(query).toHaveBeenCalledWith(STYLE_MEDIA_QUERY.PHONE);
  });

  it('is false when media queries are unavailable', () => {
    Object.defineProperty(window, 'matchMedia', { configurable: true, value: undefined });
    expect(getMobileSnapshot()).toBe(false);
  });
});
