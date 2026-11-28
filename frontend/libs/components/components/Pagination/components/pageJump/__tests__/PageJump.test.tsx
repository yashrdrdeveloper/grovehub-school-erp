import { fireEvent, render, screen } from '@testing-library/react';

import { PageJump } from '../PageJump';

describe('PageJump', () => {
  it('allows arbitrary valid page jumps and rejects invalid page numbers', () => {
    const onChange = jest.fn();
    render(<PageJump currentPage={1} totalPages={40} disabled={false} onChange={onChange} />);
    const input = screen.getByRole('spinbutton', { name: 'Page number' });
    fireEvent.change(input, { target: { value: '41' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(onChange).not.toHaveBeenCalled();
    fireEvent.change(input, { target: { value: '17' } });
    fireEvent.click(screen.getByRole('button', { name: 'Go' }));
    expect(onChange).toHaveBeenCalledWith(17);
  });
});
