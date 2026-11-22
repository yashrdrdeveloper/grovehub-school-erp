import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import type { ReactElement } from 'react';
import { Link } from '@Components/router';
import { Pin, PinOff } from '@Components/icons';
import { APP_LAUNCHER_COPY } from '../constants/appLauncher.copy';
import { getAppTestIds } from '../helpers/appLauncher.getAppTestIds';
import type { RenderAppLauncherApp } from './types/appLauncher.app';

/** Stateless catalog card. The caller supplies its pinned state explicitly. */
export const renderAppLauncherApp: RenderAppLauncherApp = ({ app, isPinned, onOpenChange, onTogglePin }): ReactElement => {
  const Icon = app.icon;
  let pinLabel = 'Pin';
  let PinIcon = Pin;
  let pinClass = 'catalog-app__pin';
  if (isPinned) {
    pinLabel = 'Unpin';
    PinIcon = PinOff;
    pinClass += ' catalog-app__pin--active';
  }
  const ids = getAppTestIds(app.path);
  const close = () => onOpenChange(false);
  const toggle = () => onTogglePin(app.path);
  return (
    <Box as="div" key={app.path} className="catalog-app" data-testid={ids.card}>
      <Link to={app.path} className="catalog-app__link" onClick={close}>
        <Box as="span" className="catalog-app__icon"><Icon size={22} strokeWidth={1.8}/></Box>
        <Box as="strong">{app.name}</Box><Box as="small">{app.description}</Box>
        <Box as="span" className="catalog-app__status">{APP_LAUNCHER_COPY.STATUS}</Box>
      </Link>
      <Button unstyled type="button" className={pinClass} onClick={toggle}
        aria-label={pinLabel + ' ' + app.name} title={pinLabel + ' from sidebar'} data-testid={ids.pin}>
        <PinIcon size={15}/>
      </Button>
    </Box>
  );
};
