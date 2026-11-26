import { fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { MobileNav } from '../MobileNav';

import type { MobileNavProps } from '../types/mobileNav';

describe('MobileNav', () => {
  it('retains Home navigation and delegates Apps/Menu events to the shell model', () => {
    const openCatalog = jest.fn();
    const openSidebar = jest.fn();
    const model = { openCatalog, openSidebar, sidebarOpen: false } as unknown as MobileNavProps['model'];

    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <MobileNav model={model} />
      </MemoryRouter>,
    );

    const navigation = screen.getByRole('navigation', { name: 'Mobile navigation' });
    expect(navigation).toBeTruthy();
    expect(screen.getByRole('link', { name: 'Home' }).getAttribute('href')).toBe('/dashboard');

    fireEvent.click(screen.getByRole('button', { name: 'Apps' }));
    const menu = screen.getByRole('button', { name: 'Menu' });
    expect(menu.getAttribute('aria-controls')).toBe('workspace-sidebar');
    expect(menu.getAttribute('aria-expanded')).toBe('false');
    fireEvent.click(menu);

    expect(openCatalog).toHaveBeenCalledTimes(1);
    expect(openSidebar).toHaveBeenCalledTimes(1);
  });
});
