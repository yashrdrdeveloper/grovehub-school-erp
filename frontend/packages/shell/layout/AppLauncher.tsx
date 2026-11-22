import { useEffect, useState } from 'react';
import * as DialogPrimitive from '@radix-ui/react-dialog';
import { Grid2X2, Pin, PinOff, Search, X } from '@Components/icons';
import { Link } from 'react-router-dom';

import { appCatalog, type CatalogApp } from '../registry/appCatalog.registry';

export { appCatalog, defaultPinnedApps } from '../registry/appCatalog.registry';
export type { CatalogApp } from '../registry/appCatalog.registry';

export interface AppLauncherProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly pinned: string[];
  readonly onTogglePin: (path: string) => void;
  readonly sidebarCollapsed: boolean;
}

/** Catalog lists implemented modules only. The shell remains mounted around its outlet. */
export const AppLauncher = ({ open, onOpenChange, pinned, onTogglePin, sidebarCollapsed }: AppLauncherProps) => {
  const [filter, setFilter] = useState('');
  const matching = appCatalog.filter(app =>
    (app.name + ' ' + app.description).toLowerCase().includes(filter.trim().toLowerCase()));
  const pinnedApps = matching.filter(app => pinned.includes(app.path));
  const otherApps = matching.filter(app => !pinned.includes(app.path));

  useEffect(() => { if (open) setFilter(''); }, [open]);

  const renderApp = (app: CatalogApp) => {
    const Icon = app.icon;
    const isPinned = pinned.includes(app.path);
    let pinLabel = 'Pin';
    let PinIcon = Pin;
    if (isPinned) {
      pinLabel = 'Unpin';
      PinIcon = PinOff;
    }
    return (
      <div key={app.path} className="catalog-app" data-testid={`catalog-app-${app.name.toLowerCase()}`}>
        <Link to={app.path} className="catalog-app__link" onClick={() => onOpenChange(false)}>
          <span className="catalog-app__icon"><Icon size={22} strokeWidth={1.8}/></span>
          <strong>{app.name}</strong><small>{app.description}</small>
          <span className="catalog-app__status">Open app</span>
        </Link>
        <button type="button" className={`catalog-app__pin${isPinned ? ' catalog-app__pin--active' : ''}`}
          onClick={() => onTogglePin(app.path)} aria-label={`${pinLabel} ${app.name}`}
          title={`${pinLabel} from sidebar`} data-testid={`catalog-pin-${app.name.toLowerCase()}`}>
          <PinIcon size={15}/>
        </button>
      </div>
    );
  };

  return (
    <DialogPrimitive.Root open={open} onOpenChange={onOpenChange}>
      <DialogPrimitive.Portal>
        <DialogPrimitive.Overlay className={`catalog-overlay${sidebarCollapsed ? ' catalog-overlay--collapsed' : ''}`}/>
        <DialogPrimitive.Content className={`catalog-content${sidebarCollapsed ? ' catalog-content--collapsed' : ''}`} aria-describedby="catalog-description">
          <div className="catalog-content__header">
            <div className="catalog-content__title">
              <span className="catalog-content__mark"><Grid2X2 size={23}/></span>
              <div><DialogPrimitive.Title>App catalog</DialogPrimitive.Title>
                <DialogPrimitive.Description id="catalog-description">Open your available RDR School applications and pin them for quick access.</DialogPrimitive.Description></div>
            </div>
            <DialogPrimitive.Close className="icon-button catalog-content__close" aria-label="Close app catalog"><X size={20}/></DialogPrimitive.Close>
          </div>
          <div className="catalog-content__body">
            <div className="catalog-content__tools">
              <div className="catalog-search"><Search size={18}/><input value={filter} onChange={event => setFilter(event.target.value)}
                placeholder="Search applications..." aria-label="Search applications" data-testid="catalog-search"/></div>
              <span className="catalog-content__count">{appCatalog.length} available</span>
            </div>
            {pinnedApps.length > 0 && (
              <section className="catalog-section">
                <div className="catalog-section__heading"><div><Pin size={16}/><h3>Pinned apps</h3></div><span>Quick access in your sidebar</span></div>
                <div className="catalog-app-grid">{pinnedApps.map(renderApp)}</div>
              </section>
            )}
            {otherApps.length > 0 && (
              <section className="catalog-section">
                <div className="catalog-section__heading"><div><Grid2X2 size={16}/><h3>Applications</h3></div><span>Available in your school workspace</span></div>
                <div className="catalog-app-grid">{otherApps.map(renderApp)}</div>
              </section>
            )}
            {matching.length === 0 && <p className="catalog-empty">No applications match your search.</p>}
          </div>
        </DialogPrimitive.Content>
      </DialogPrimitive.Portal>
    </DialogPrimitive.Root>
  );
};
