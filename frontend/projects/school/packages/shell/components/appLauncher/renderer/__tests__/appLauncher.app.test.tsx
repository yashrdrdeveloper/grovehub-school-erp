import { fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { appCatalog } from '../../../../constants/shell.catalog';
import { renderAppLauncherApp } from '../appLauncher.app';

describe('renderAppLauncherApp', () => {
  it('renders a catalog link and its independent pin action', () => {
    const interactions = { toggled: '', open: true };
    const app = appCatalog[0];
    if (!app) throw new Error('Catalog must contain the dashboard');

    render(
      <MemoryRouter>
        {renderAppLauncherApp({
          app,
          isPinned: false,
          onOpenChange: (next) => {
            interactions.open = next;
          },
          onTogglePin: (path) => {
            interactions.toggled = path;
          },
        })}
      </MemoryRouter>,
    );

    expect(screen.getByRole('link', { name: /Dashboard/ }).getAttribute('href')).toBe(app.path);
    fireEvent.click(screen.getByRole('button', { name: 'Pin Dashboard' }));
    expect(interactions.toggled).toBe(app.path);
    fireEvent.click(screen.getByRole('link', { name: /Dashboard/ }));
    expect(interactions.open).toBe(false);
  });

  it('exposes the unpin action for pinned applications', () => {
    const app = appCatalog[0];
    if (!app) throw new Error('Catalog must contain the dashboard');

    render(
      <MemoryRouter>
        {renderAppLauncherApp({
          app,
          isPinned: true,
          onOpenChange: () => {},
          onTogglePin: () => {},
        })}
      </MemoryRouter>,
    );

    expect(screen.getByRole('button', { name: 'Unpin Dashboard' })).toBeTruthy();
  });
});
