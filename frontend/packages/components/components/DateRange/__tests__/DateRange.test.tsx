import { fireEvent, render, screen, within } from '@testing-library/react';

import { DateRange } from '../DateRange';

describe('DateRange', () => {
  it('displays the optional header and footer independently', () => {
    const onChange = jest.fn();
    render(<DateRange value={{ startDate: '2026-11-01', endDate: '2026-11-30', presetId: 'thisMonth' }} onChange={onChange} />);

    fireEvent.click(screen.getByRole('button', { name: 'Select date range' }));
    const dialog = screen.getByRole('dialog', { name: 'Date range' });
    const header = within(dialog).getByRole('button', { name: 'Close' }).closest('[data-custom-heading="true"]') as HTMLElement;
    const apply = within(dialog).getByRole('button', { name: 'Apply' });

    expect(within(header).getByText('01 Nov 2026 – 30 Nov 2026')).toBeTruthy();
    expect(within(header).queryByRole('button', { name: 'Apply' })).toBeNull();
    expect(within(dialog).getByRole('button', { name: 'Clear' })).toBeTruthy();
    fireEvent.click(apply);
    expect(onChange).toHaveBeenCalledWith({ startDate: '2026-11-01', endDate: '2026-11-30', presetId: 'thisMonth' });
  });

  it('omits both sections and immediately commits a preset in filter mode', () => {
    const onChange = jest.fn();
    render(<DateRange showHeader={false} showFooter={false} value={{ startDate: null, endDate: null }} onChange={onChange} />);

    fireEvent.click(screen.getByRole('button', { name: 'Select date range' }));
    const dialog = screen.getByRole('dialog', { name: 'Date range' });
    expect(within(dialog).queryByRole('button', { name: 'Close' })).toBeNull();
    expect(within(dialog).queryByRole('button', { name: 'Clear' })).toBeNull();
    expect(within(dialog).queryByRole('button', { name: 'Apply' })).toBeNull();

    fireEvent.click(within(dialog).getByRole('button', { name: 'Last 30 days' }));
    expect(onChange).toHaveBeenCalledTimes(1);
    expect(onChange).toHaveBeenCalledWith(expect.objectContaining({ presetId: 'last30Days' }));
    expect(screen.queryByRole('dialog', { name: 'Date range' })).toBeNull();
  });
});
