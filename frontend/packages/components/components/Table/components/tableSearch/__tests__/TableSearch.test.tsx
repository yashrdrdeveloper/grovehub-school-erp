import { fireEvent, render, screen } from '@testing-library/react';

import { TableSearch } from '../TableSearch';

describe('TableSearch', () => {
  it('shows searchable fields as a hover tooltip without opening a dialog on click', () => {
    render(<TableSearch value="" placeholder="Search" debounceMs={300} fields={['username', 'email']} onChange={jest.fn()} />);

    const info = screen.getByRole('button', { name: 'Search is based on username, email.' });
    expect(screen.queryByRole('tooltip')).toBeNull();

    fireEvent.mouseEnter(info);
    expect(screen.getByRole('tooltip').textContent).toBe('Search is based on username, email.');
    fireEvent.pointerDown(info);
    fireEvent.click(info);
    expect(screen.queryByRole('tooltip')).toBeNull();
    expect(screen.queryByRole('dialog')).toBeNull();
  });
});
