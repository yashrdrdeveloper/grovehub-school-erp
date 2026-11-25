import { fireEvent, render, screen, within } from '@testing-library/react';

import { DateRange } from '../DateRange';

describe('DateRange', () => {
  it('displays the formatted range and actions in the header while retaining the accessible dialog title', () => {
    const onChange = jest.fn();
    render(<DateRange value={{ startDate: '2026-11-01', endDate: '2026-11-30', presetId: 'thisMonth' }} onChange={onChange} />);

    fireEvent.click(screen.getByRole('button', { name: 'Select date range' }));
    const dialog = screen.getByRole('dialog', { name: 'Date range' });
    const apply = within(dialog).getByRole('button', { name: 'Apply' });
    const header = apply.closest('[data-custom-heading="true"]') as HTMLElement;

    expect(header).toBeTruthy();
    expect(within(header).getByText('01 Nov 2026 – 30 Nov 2026')).toBeTruthy();
    expect(within(header).getByRole('button', { name: 'Clear' })).toBeTruthy();
    expect(within(header).getByRole('button', { name: 'Close' })).toBeTruthy();
    expect(within(dialog).getAllByRole('button', { name: 'Apply' })).toHaveLength(1);

    fireEvent.click(apply);
    expect(onChange).toHaveBeenCalledWith({ startDate: '2026-11-01', endDate: '2026-11-30', presetId: 'thisMonth' });
  });
});
