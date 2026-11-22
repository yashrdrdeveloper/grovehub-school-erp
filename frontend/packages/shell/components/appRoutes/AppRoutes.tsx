import { APP_ROUTE } from '@Business/app/routes';
import { Navigate, Route, Routes } from 'react-router-dom';

import { PublicOnly, RequireSession, useRecoveryFlowBoundary } from '@Domain/auth';
import { LandingPage, PUBLIC_SITE_PATHS } from '@Domain/publicSite';

import { publicRoutes, workspaceRoutes } from '../../registry';
import { ProtectedLayout } from '../protectedLayout';
import type { AppRoutesComponent } from './types/appRoutes';

/** Domain roots publish routes; shell only composes public/protected envelopes. */
export const AppRoutes: AppRoutesComponent = () => {
  useRecoveryFlowBoundary();
  const landing = PUBLIC_SITE_PATHS.map(path => <Route key={path} path={path} element={<LandingPage/>}/>);
  const publicPages = publicRoutes.map(({ path, element }) =>
    <Route key={path} path={path} element={element}/>);
  const workspacePages = workspaceRoutes.map(({ path, element }) =>
    <Route key={path} path={path} element={element}/>);

  return (
    <Routes>
      {landing}
      <Route element={<PublicOnly/>}>{publicPages}</Route>
      <Route element={<RequireSession/>}>
        <Route element={<ProtectedLayout/>}>{workspacePages}</Route>
      </Route>
      <Route path="*" element={<Navigate to={APP_ROUTE.HOME} replace/>}/>
    </Routes>
  );
};
