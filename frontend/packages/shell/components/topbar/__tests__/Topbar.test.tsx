import { fireEvent, render, screen } from '@testing-library/react';

import { Topbar } from '../Topbar';

import type { TopbarProps } from '../types/topbar';

describe('Topbar', () => {
  it('delegates menu and account actions to the shell model', () => {
    const openSidebar = jest.fn();
    const openAccountPreferences = jest.fn();
    const model = { openSidebar, openAccountPreferences } as unknown as TopbarProps['model'];

    render(<Topbar model={model} schoolLabel="RDR School" accountLabel="parent@example.test" />);

    expect(screen.getByTestId('workspace-topbar')).toBeTruthy();
    expect(screen.getByText('RDR School')).toBeTruthy();

    fireEvent.click(screen.getByRole('button', { name: 'Open menu' }));
    fireEvent.click(screen.getByRole('button', { name: 'Open account preferences' }));

    expect(openSidebar).toHaveBeenCalledTimes(1);
    expect(openAccountPreferences).toHaveBeenCalledTimes(1);
  });
});
