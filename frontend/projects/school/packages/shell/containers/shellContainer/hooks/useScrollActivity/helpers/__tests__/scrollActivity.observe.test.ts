import { observeScrollActivity } from '../scrollActivity.observe';

describe('observeScrollActivity', () => {
  beforeEach(() => jest.useFakeTimers());
  afterEach(() => jest.useRealTimers());

  it('reveals nested scrollbars independently and resets the idle timeout', () => {
    const first = document.createElement('div');
    const second = document.createElement('div');
    document.body.append(first, second);
    const dispose = observeScrollActivity(document);
    try {
      first.dispatchEvent(new Event('scroll'));
      expect(first.getAttribute('data-gh-scroll-active')).toBe('true');
      expect(second.hasAttribute('data-gh-scroll-active')).toBe(false);
      jest.advanceTimersByTime(800);
      first.dispatchEvent(new Event('scroll'));
      second.dispatchEvent(new Event('scroll'));
      jest.advanceTimersByTime(800);
      expect(first.hasAttribute('data-gh-scroll-active')).toBe(true);
      expect(second.hasAttribute('data-gh-scroll-active')).toBe(true);
      jest.advanceTimersByTime(200);
      expect(first.hasAttribute('data-gh-scroll-active')).toBe(false);
      expect(second.hasAttribute('data-gh-scroll-active')).toBe(false);
    } finally {
      dispose();
      first.remove();
      second.remove();
    }
  });

  it('handles document scrolling and removes listeners and active state on cleanup', () => {
    const dispose = observeScrollActivity(document);
    document.dispatchEvent(new Event('scroll'));
    expect(document.documentElement.hasAttribute('data-gh-scroll-active')).toBe(true);
    dispose();
    expect(document.documentElement.hasAttribute('data-gh-scroll-active')).toBe(false);
    document.dispatchEvent(new Event('scroll'));
    expect(document.documentElement.hasAttribute('data-gh-scroll-active')).toBe(false);
    expect(jest.getTimerCount()).toBe(0);
  });
});
