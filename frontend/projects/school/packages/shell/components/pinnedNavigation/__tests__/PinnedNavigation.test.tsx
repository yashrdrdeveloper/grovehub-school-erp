import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { PinnedNavigation } from '../PinnedNavigation';

describe('PinnedNavigation', () => {
  it('shows the accessible empty state without a phantom application', () => {
    render(
      <MemoryRouter>
        <PinnedNavigation pinned={[]} closeSidebar={jest.fn()} togglePin={jest.fn()} />
      </MemoryRouter>,
    );

    expect(screen.getByRole('navigation', { name: 'Pinned applications' })).toBeTruthy();
    expect(screen.getByText('Pin apps from the catalog to see them here.')).toBeTruthy();
    expect(screen.queryByRole('link', { name: 'Dashboard' })).toBeNull();
  });

  it('uses the implemented route metadata for pinned app navigation', () => {
    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <PinnedNavigation pinned={['/dashboard']} closeSidebar={jest.fn()} togglePin={jest.fn()} />
      </MemoryRouter>,
    );

    expect(screen.getByRole('link', { name: 'Dashboard' }).getAttribute('href')).toBe('/dashboard');
    expect(screen.getByRole('button', { name: 'Unpin Dashboard' })).toBeTruthy();
  });
});
