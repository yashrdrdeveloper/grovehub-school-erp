import { fireEvent, render, screen } from '@testing-library/react';

import { Topbar } from '../Topbar';

import type { TopbarProps } from '../types/topbar';

describe('Topbar', () => {
  it('displays the school monogram without a breadcrumb or account control and retains mobile navigation', () => {
    const openSidebar = jest.fn();
    const model = { openSidebar } as unknown as TopbarProps['model'];

    render(<Topbar model={model} schoolLabel="RDR School" />);

    expect(screen.getByTestId('workspace-topbar')).toBeTruthy();
    expect(screen.getByRole('img', { name: 'RDR School' })).toBeTruthy();
    expect(screen.queryByText('RDR School')).toBeNull();
    expect(screen.queryByText('Workspace')).toBeNull();
    expect(screen.queryByRole('button', { name: 'Open account preferences' })).toBeNull();

    fireEvent.click(screen.getByRole('button', { name: 'Open menu' }));
    expect(openSidebar).toHaveBeenCalledTimes(1);
  });
});
