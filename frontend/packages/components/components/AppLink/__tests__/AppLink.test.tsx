import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { AppLink } from '../AppLink';

describe('AppLink', () => {
  it('retains native anchor navigation and forwards the selector', () => {
    render(<MemoryRouter><AppLink to="/dashboard" dataTestId="dashboard-link">Dashboard</AppLink></MemoryRouter>);
    expect(screen.getByRole('link',{name:'Dashboard'}).getAttribute('href')).toBe('/dashboard');
    expect(screen.getByTestId('dashboard-link').tagName).toBe('A');
  });
});
