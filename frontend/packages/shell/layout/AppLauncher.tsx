import { useEffect, useRef, useState } from 'react';
import * as DialogPrimitive from '@radix-ui/react-dialog';
import { Grid2X2, Pin, PinOff, Search, X } from 'lucide-react';
import { Link } from 'react-router-dom';

import { appCatalog, type CatalogApp } from '../registry/appCatalog.registry';

// Compatibility exports for existing launcher consumers.
export { appCatalog, defaultPinnedApps } from '../registry/appCatalog.registry';
export type { CatalogApp } from '../registry/appCatalog.registry';

const categories = ['All', 'Workspace', 'Student life', 'Operations', 'Administration'];

export interface AppLauncherProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  pinned: string[];
  onTogglePin: (path: string) => void;
  sidebarCollapsed: boolean;
}

export const AppLauncher = ({ open, onOpenChange, pinned, onTogglePin, sidebarCollapsed }: AppLauncherProps) => {
  const [filter, setFilter] = useState('');
  const [category, setCategory] = useState('All');
  const searchRef = useRef<HTMLInputElement>(null);
  const matching = appCatalog.filter(app =>
    (category === 'All' || app.category === category) &&
    (app.name + ' ' + app.description).toLowerCase().includes(filter.trim().toLowerCase()),
  );
  const pinnedApps = matching.filter(app => app.path && pinned.includes(app.path));
  const otherApps = matching.filter(app => !app.path || !pinned.includes(app.path));
  const available = appCatalog.filter(app => app.path).length;

  useEffect(() => {
    if (open) setFilter('');
  }, [open]);

  const renderApp = (app: CatalogApp) => {
    const Icon = app.icon;
    const isPinned = !!app.path && pinned.includes(app.path);
    return (
      <div key={app.name} className={`catalog-app${app.path ? '' : ' catalog-app--planned'}`}>
        {app.path ? (
          <Link to={app.path} className="catalog-app__link" onClick={() => onOpenChange(false)}>
            <span className="catalog-app__icon"><Icon size={22} strokeWidth={1.8}/></span>
            <strong>{app.name}</strong>
            <small>{app.description}</small>
            <span className="catalog-app__status">Open app</span>
          </Link>
        ) : (
          <div className="catalog-app__link" aria-label={`${app.name} planned`}>
            <span className="catalog-app__icon"><Icon size={22} strokeWidth={1.8}/></span>
            <strong>{app.name}</strong>
            <small>{app.description}</small>
            <span className="catalog-app__status">Coming soon</span>
          </div>
        )}
        {app.path && (
          <button
            type="button"
            className={`catalog-app__pin${isPinned ? ' catalog-app__pin--active' : ''}`}
            onClick={() => onTogglePin(app.path!)}
            aria-label={`${isPinned ? 'Unpin' : 'Pin'} ${app.name}`}
            title={isPinned ? 'Unpin from sidebar' : 'Pin to sidebar'}
          >{isPinned ? <PinOff size={15}/> : <Pin size={15}/>}</button>
        )}
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
                <DialogPrimitive.Description id="catalog-description">Everything at RDR School, in one place. Pin the apps you use most.</DialogPrimitive.Description>
              </div>
            </div>
            <DialogPrimitive.Close className="icon-button catalog-content__close" aria-label="Close app catalog"><X size={20}/></DialogPrimitive.Close>
          </div>
          <div className="catalog-content__body">
            <div className="catalog-content__tools">
              <div className="catalog-search"><Search size={18}/><input ref={searchRef} value={filter} onChange={event => setFilter(event.target.value)} placeholder="Search school applications..." aria-label="Search applications"/></div>
              <span className="catalog-content__count">{available} available · {appCatalog.length - available} coming soon</span>
            </div>
            <div className="catalog-tabs" role="group" aria-label="Filter by category">
              {categories.map(item => <button key={item} type="button" className={category === item ? 'catalog-tab catalog-tab--active' : 'catalog-tab'} aria-pressed={category === item} onClick={() => setCategory(item)}>{item}</button>)}
            </div>
            {pinnedApps.length > 0 && (
              <section className="catalog-section">
                <div className="catalog-section__heading"><div><Pin size={16}/><h3>Pinned apps</h3></div><span>Quick access in your sidebar</span></div>
                <div className="catalog-app-grid">{pinnedApps.map(renderApp)}</div>
              </section>
            )}
            {otherApps.length > 0 && (
              <section className="catalog-section">
                <div className="catalog-section__heading"><div><Grid2X2 size={16}/><h3>{pinnedApps.length > 0 ? 'Explore more apps' : 'All applications'}</h3></div><span>Browse your school workspace</span></div>
                <div className="catalog-app-grid">{otherApps.map(renderApp)}</div>
              </section>
            )}
            {matching.length === 0 && <p className="catalog-empty">No applications match your search. Try another name or category.</p>}
          </div>
        </DialogPrimitive.Content>
      </DialogPrimitive.Portal>
    </DialogPrimitive.Root>
  );
};
