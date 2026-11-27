import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

import { AuthProvider } from '../../providers/authProvider';

import { PublicOnly, RequireSession } from '../RequireSession';

const me = jest.fn<Promise<unknown>, []>();
const refresh = jest.fn<Promise<unknown>, []>();
jest.mock('@Business/auth', () => ({
  authenticationApi: {
    me: () => me(),
    refresh: () => refresh(),
  },
}));

const renderAt = (path: string) =>
  render(
    <MemoryRouter initialEntries={[path]}>
      <AuthProvider>
        <Routes>
          <Route element={<PublicOnly />}>
            <Route path="/set-password" element={<div>Set password page</div>} />
            <Route path="/login" element={<div>Login page</div>} />
          </Route>
          <Route element={<RequireSession />}>
            <Route path="/dashboard" element={<div>Private dashboard</div>} />
          </Route>
        </Routes>
      </AuthProvider>
    </MemoryRouter>,
  );

describe('RequireSession', () => {
  beforeEach(() => {
    me.mockReset();
    refresh.mockReset();
    me.mockRejectedValue(new Error('No session'));
    refresh.mockRejectedValue(new Error('No refresh cookie'));
  });

  it('does not request an access or refresh session on the set-password page', async () => {
    renderAt('/set-password');
    expect(screen.getByText('Set password page')).toBeTruthy();
    expect(me).not.toHaveBeenCalled();
    expect(refresh).not.toHaveBeenCalled();
  });

  it('only checks access and refresh sessions when entering a protected route', async () => {
    renderAt('/dashboard');
    expect(await screen.findByText('Login page')).toBeTruthy();
    expect(me).toHaveBeenCalledTimes(1);
    expect(refresh).toHaveBeenCalledTimes(1);
  });

  it('renders the protected route when a session exists', async () => {
    me.mockResolvedValue({ email: 'student@school.edu' });
    renderAt('/dashboard');
    expect(await screen.findByText('Private dashboard')).toBeTruthy();
    await waitFor(() => expect(me).toHaveBeenCalledTimes(1));
    expect(refresh).not.toHaveBeenCalled();
  });
});
