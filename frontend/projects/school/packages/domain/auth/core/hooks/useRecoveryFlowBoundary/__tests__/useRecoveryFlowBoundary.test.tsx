import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Link, MemoryRouter, Route, Routes } from 'react-router-dom';

import { useRecoveryFlowBoundary } from '../useRecoveryFlowBoundary';
import {
  clearRecoveryFlow,
  readRecoveryFlow,
  saveRecoveryProof,
} from '../../../session/recoveryFlowSession';

const Boundary = () => {
  useRecoveryFlowBoundary();
  return null;
};

const renderRoutes = (path: string) =>
  render(
    <MemoryRouter initialEntries={[path]}>
      <Boundary />
      <Routes>
        <Route
          path="/set-password"
          element={
            <>
              <p>Setup</p>
              <Link to="/reset-password">Reset</Link>
              <Link to="/login">Login</Link>
            </>
          }
        />
        <Route
          path="/reset-password"
          element={
            <>
              <p>Reset</p>
              <Link to="/set-password">Setup</Link>
            </>
          }
        />
        <Route path="/login" element={<p>Login</p>} />
      </Routes>
    </MemoryRouter>,
  );

describe('useRecoveryFlowBoundary', () => {
  beforeEach(() => {
    clearRecoveryFlow();
    window.sessionStorage.clear();
  });

  it('retains its flow when the same page is loaded again', () => {
    saveRecoveryProof('invitation', 'private-proof');
    renderRoutes('/set-password');
    expect(readRecoveryFlow('invitation')).toMatchObject({
      stage: 'proof',
      proof: 'private-proof',
    });
    expect(window.sessionStorage.length).toBe(0);
  });

  it('clears the previous proof when switching set-password to reset-password', async () => {
    saveRecoveryProof('invitation', 'private-proof');
    renderRoutes('/set-password');
    fireEvent.click(screen.getByText('Reset'));
    await waitFor(() => expect(readRecoveryFlow('recovery')).toBeNull());
    expect(window.sessionStorage.length).toBe(0);
  });

  it('clears a proof when navigating to login', async () => {
    saveRecoveryProof('invitation', 'private-proof');
    renderRoutes('/set-password');
    fireEvent.click(screen.getByText('Login'));
    await waitFor(() => expect(readRecoveryFlow('invitation')).toBeNull());
  });

  it('does not restore a different-purpose proof on direct navigation', () => {
    saveRecoveryProof('recovery', 'reset-private-proof');
    renderRoutes('/set-password');
    expect(readRecoveryFlow('invitation')).toBeNull();
  });
});
