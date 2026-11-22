import { fireEvent, render, screen } from '@testing-library/react';
import { Link, MemoryRouter, Outlet, Route, Routes } from 'react-router-dom';
import { ProtectedLayout } from '../ProtectedLayout';

jest.mock('@Domain/auth', () => ({
  useAuth: () => ({ account: { email: 'school@example.test' }, signOut: jest.fn() }),
}));
jest.mock('../../../containers/workspaceShellContainer', () => ({
  WorkspaceShellContainer: () => <div data-testid="persistent-workspace-shell"><Outlet/></div>,
}));

describe('ProtectedLayout', () => {
  it('retains the same mounted shell while protected child routes change', () => {
    render(<MemoryRouter initialEntries={['/dashboard']}>
      <Routes>
        <Route element={<ProtectedLayout/>}>
          <Route path="/dashboard" element={<><h1>Dashboard</h1><Link to="/workspace-test">Next app</Link></>}/>
          <Route path="/workspace-test" element={<h1>Second app</h1>}/>
        </Route>
      </Routes>
    </MemoryRouter>);
    const root = screen.getByTestId('persistent-workspace-shell');
    fireEvent.click(screen.getByRole('link', {name:'Next app'}));
    expect(screen.getByRole('heading',{name:'Second app'})).toBeTruthy();
    expect(screen.getByTestId('persistent-workspace-shell')).toBe(root);
  });
});
