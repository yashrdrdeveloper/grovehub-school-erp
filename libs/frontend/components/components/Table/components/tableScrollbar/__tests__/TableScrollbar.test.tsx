import { fireEvent, render, screen } from '@testing-library/react';

import { TableScrollbar } from '../TableScrollbar';

describe('TableScrollbar', () => {
  it('does not render when the table has no overflow', () => {
    const element = document.createElement('div');
    Object.defineProperties(element, {
      clientHeight: { value: 200 },
      scrollHeight: { value: 200 },
    });
    render(<TableScrollbar viewport={{ current: element }} viewportId="table" axis="vertical" />);
    expect(screen.queryByRole('scrollbar')).toBeNull();
  });

  it('uses native scroll offsets and supports keyboard navigation', () => {
    const element = document.createElement('div');
    Object.defineProperties(element, {
      clientHeight: { value: 200 },
      scrollHeight: { value: 400 },
    });
    render(<TableScrollbar viewport={{ current: element }} viewportId="table" axis="vertical" />);
    const scrollbar = screen.getByRole('scrollbar', { name: 'Vertical table scrollbar' });
    fireEvent.keyDown(scrollbar, { key: 'End' });
    expect(element.scrollTop).toBe(200);
    fireEvent.scroll(element);
    expect(scrollbar.getAttribute('aria-valuenow')).toBe('200');
    fireEvent.keyDown(scrollbar, { key: 'Home' });
    expect(element.scrollTop).toBe(0);
  });
});
