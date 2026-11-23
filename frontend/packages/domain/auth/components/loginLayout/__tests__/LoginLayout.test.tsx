import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { LoginLayout } from '../LoginLayout';

describe('LoginLayout', () => {
  it('keeps one semantic heading with a distinct mobile title and fallback', () => {
    const screenView = render(
      <MemoryRouter>
        <LoginLayout
          title="Set a new password."
          mobileTitle="Set password"
          description="Choose a strong password."
        >
          <span>Form</span>
        </LoginLayout>
      </MemoryRouter>,
    );

    const heading = screen.getByRole('heading', { level: 1 });
    expect(heading.id).toBe('rdr-auth-title');
    expect(heading.querySelectorAll('span')).toHaveLength(2);
    expect(heading.textContent).toContain('Set a new password.');
    expect(heading.textContent).toContain('Set password');

    screenView.rerender(
      <MemoryRouter>
        <LoginLayout title="Welcome back." description="Sign in.">
          <span>Form</span>
        </LoginLayout>
      </MemoryRouter>,
    );
    expect(screen.getByRole('heading', { level: 1 }).querySelectorAll('span')).toHaveLength(2);
    expect(screen.getByRole('heading', { level: 1 }).textContent).toBe('Welcome back.Welcome back.');
  });
});
