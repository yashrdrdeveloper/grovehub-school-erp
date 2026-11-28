import { fireEvent, render, screen } from '@testing-library/react';

import { ThemePreference } from '@Components/hooks/useWorkspaceTheme';

import { Topbar } from '../Topbar';

import type { TopbarProps } from '../types/topbar';

describe('Topbar', () => {
  it('displays the school brand and theme choices without restoring the account control', () => {
    const openSidebar = jest.fn();
    const setPreference = jest.fn();
    const model = {
      openSidebar,
      preference: ThemePreference.SYSTEM,
      setPreference,
    } as unknown as TopbarProps['model'];

    render(<Topbar model={model} schoolLabel="RDR School" />);

    expect(screen.getByTestId('workspace-topbar')).toBeTruthy();
    expect(screen.getByRole('img', { name: 'RDR School' })).toBeTruthy();
    expect(screen.getByText('School')).toBeTruthy();
    expect(screen.queryByText('Workspace')).toBeNull();
    expect(screen.queryByRole('button', { name: 'Open account preferences' })).toBeNull();

    expect(screen.getByRole('button', { name: 'System' }).getAttribute('aria-pressed')).toBe(
      'true',
    );
    fireEvent.click(screen.getByRole('button', { name: 'Light' }));
    fireEvent.click(screen.getByRole('button', { name: 'Dark' }));
    fireEvent.click(screen.getByRole('button', { name: 'System' }));
    expect(setPreference.mock.calls).toEqual([
      [ThemePreference.LIGHT],
      [ThemePreference.DARK],
      [ThemePreference.SYSTEM],
    ]);

    fireEvent.click(screen.getByRole('button', { name: 'Open menu' }));
    expect(openSidebar).toHaveBeenCalledTimes(1);
  });
});
