import { fireEvent, render, screen } from '@testing-library/react';

import { Dialog } from '../Dialog';

describe('Dialog', () => {
  it('renders a titled accessible portal with a working close action', () => {
    const onOpenChange = jest.fn();

    render(
      <Dialog
        open
        onOpenChange={onOpenChange}
        title="Applications"
        description="Manage the catalog"
        dataTestId="catalog-dialog"
      >
        <p>Choose an app</p>
      </Dialog>,
    );

    expect(screen.getByRole('dialog', { name: 'Applications' })).toBeTruthy();
    expect(screen.getByTestId('catalog-dialog')).toBeTruthy();

    fireEvent.click(screen.getByRole('button', { name: 'Close' }));
    expect(onOpenChange).toHaveBeenCalledWith(false);
  });
});
