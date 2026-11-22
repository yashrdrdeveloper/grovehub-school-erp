import { act, render } from '@testing-library/react';

import { useAuthViewportFit } from '../useAuthViewportFit';

const Harness = () => {
  const { contentRef, toolbarRef, mainRef } = useAuthViewportFit();
  return <section ref={contentRef} data-testid="auth-panel">
    <div ref={toolbarRef} data-testid="auth-toolbar">Back to school</div>
    <div ref={mainRef} data-testid="auth-main">Welcome back.</div>
  </section>;
};

describe('useAuthViewportFit', () => {
  const originalMatchMedia = window.matchMedia;
  const originalAnimationFrame = window.requestAnimationFrame;
  const originalCancelFrame = window.cancelAnimationFrame;

  beforeEach(() => {
    Object.defineProperty(window, 'matchMedia', {
      configurable: true,
      value: (query: string) => ({
        matches: query === '(min-width: 721px)',
        media: query,
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
      }),
    });
    Object.defineProperty(window, 'requestAnimationFrame', {
      configurable: true,
      value: (callback: FrameRequestCallback) => { callback(0); return 1; },
    });
    Object.defineProperty(window, 'cancelAnimationFrame', {
      configurable: true,
      value: jest.fn(),
    });
  });

  afterAll(() => {
    Object.defineProperty(window, 'matchMedia', { configurable: true, value: originalMatchMedia });
    Object.defineProperty(window, 'requestAnimationFrame', { configurable: true, value: originalAnimationFrame });
    Object.defineProperty(window, 'cancelAnimationFrame', { configurable: true, value: originalCancelFrame });
  });

  it('fits a tablet form within the available panel without scrolling', () => {
    const { getByTestId } = render(<Harness/>);
    const panel = getByTestId('auth-panel');
    const toolbar = getByTestId('auth-toolbar');
    const main = getByTestId('auth-main');
    Object.defineProperty(panel, 'clientHeight', { configurable: true, get: () => 750 });
    Object.defineProperty(toolbar, 'offsetTop', { configurable: true, get: () => 20 });
    Object.defineProperty(toolbar, 'offsetHeight', { configurable: true, get: () => 48 });
    Object.defineProperty(main, 'offsetHeight', { configurable: true, get: () => 710 });

    act(() => window.dispatchEvent(new Event('resize')));

    expect(panel.getAttribute('data-auth-fit')).toBe('fit');
    const scale = Number(panel.style.getPropertyValue('--auth-panel-scale'));
    expect(scale).toBeGreaterThanOrEqual(0.82);
    expect(scale).toBeLessThan(1);
    expect(panel.style.getPropertyValue('--auth-panel-center')).toContain('px');
  });

  it('uses accessible scrolling rather than cutting off unusually tall states', () => {
    const { getByTestId } = render(<Harness/>);
    const panel = getByTestId('auth-panel');
    const toolbar = getByTestId('auth-toolbar');
    const main = getByTestId('auth-main');
    Object.defineProperty(panel, 'clientHeight', { configurable: true, get: () => 400 });
    Object.defineProperty(toolbar, 'offsetTop', { configurable: true, get: () => 20 });
    Object.defineProperty(toolbar, 'offsetHeight', { configurable: true, get: () => 48 });
    Object.defineProperty(main, 'offsetHeight', { configurable: true, get: () => 700 });

    act(() => window.dispatchEvent(new Event('resize')));

    expect(panel.getAttribute('data-auth-fit')).toBe('scroll');
    expect(panel.style.getPropertyValue('--auth-panel-scale')).toBe('');
  });
});
