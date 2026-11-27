import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { AppNavLink } from '../AppNavLink';

describe('AppNavLink', () => {
  it('forwards active navigation semantics and stable test selectors', () => {
    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <AppNavLink to="/dashboard" dataTestId="nav-dashboard">
          Dashboard
        </AppNavLink>
      </MemoryRouter>,
    );
    const link = screen.getByRole('link', { name: 'Dashboard' });
    expect(link.getAttribute('aria-current')).toBe('page');
    expect(screen.getByTestId('nav-dashboard')).toBe(link);
  });
});
