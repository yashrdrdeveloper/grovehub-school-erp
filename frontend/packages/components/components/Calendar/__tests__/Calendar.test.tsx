import { fireEvent, render, screen, within } from '@testing-library/react';

import { Calendar } from '../Calendar';

describe('Calendar', () => {
  it('navigates from month to year selection and back to day selection without committing a date', () => {
    const onChange = jest.fn();
    render(<Calendar value="2026-11-05" onChange={onChange} compact hideFooter />);
    fireEvent.click(screen.getByRole('button', { name: 'Choose month' }));
    expect(screen.getByRole('group', { name: 'Select month' })).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: 'Choose year' }));
    fireEvent.click(within(screen.getByRole('group', { name: 'Select year' })).getByRole('button', { name: '2024' }));
    fireEvent.click(within(screen.getByRole('group', { name: 'Select month' })).getByRole('button', { name: 'February 2024' }));
    expect(screen.getByRole('group', { name: 'February 2024' })).toBeTruthy();
    expect(onChange).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Thursday, February 29, 2024' }));
    expect(onChange).toHaveBeenCalledWith('2024-02-29');
  });
});
