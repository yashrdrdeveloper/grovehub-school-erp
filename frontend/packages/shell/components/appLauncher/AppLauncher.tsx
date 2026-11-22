import {
  DialogRoot, DialogPortal, DialogOverlay, DialogContent,
  DialogTitle, DialogDescription, DialogClose,
} from '@Components/components/Dialog/parts';
import { InputControl } from '@Components/components/InputControl';
import { Grid2X2, Pin, Search, X } from '@Components/icons';
import { appCatalog } from '../../registry/appCatalog.registry';
import { APP_LAUNCHER_COPY } from './constants/appLauncher.copy';
import { APP_LAUNCHER_TEST_IDS } from './constants/appLauncher.testIds';
import { renderAppLauncherApp } from './renderer/appLauncher.app';
import type { AppLauncherViewProps } from './types/appLauncher';

/** Pure catalog presentation; query state and filtering live in the container hook. */
export const AppLauncher = ({ model, open, onOpenChange, onTogglePin, sidebarCollapsed }: AppLauncherViewProps) => {
  let overlayClass = 'catalog-overlay';
  let contentClass = 'catalog-content';
  if (sidebarCollapsed) {
    overlayClass += ' catalog-overlay--collapsed';
    contentClass += ' catalog-content--collapsed';
  }
  const { matching, pinnedApps, otherApps, filter, onFilterChange } = model;
  return (
    <DialogRoot open={open} onOpenChange={onOpenChange}>
      <DialogPortal>
        <DialogOverlay className={overlayClass} data-testid={APP_LAUNCHER_TEST_IDS.OVERLAY}/>
        <DialogContent className={contentClass} aria-describedby="catalog-description" data-testid={APP_LAUNCHER_TEST_IDS.ROOT}>
          <div className="catalog-content__header" data-testid={APP_LAUNCHER_TEST_IDS.HEADER}>
            <div className="catalog-content__title">
              <span className="catalog-content__mark"><Grid2X2 size={23}/></span>
              <div><DialogTitle>{APP_LAUNCHER_COPY.TITLE}</DialogTitle>
                <DialogDescription id="catalog-description">{APP_LAUNCHER_COPY.DESCRIPTION}</DialogDescription></div>
            </div>
            <DialogClose className="icon-button catalog-content__close" aria-label="Close app catalog"><X size={20}/></DialogClose>
          </div>
          <div className="catalog-content__body" data-testid={APP_LAUNCHER_TEST_IDS.BODY}>
            <div className="catalog-content__tools">
              <div className="catalog-search"><Search size={18}/>
                <InputControl unstyled value={filter} onChange={onFilterChange}
                  placeholder={APP_LAUNCHER_COPY.SEARCH_PLACEHOLDER} aria-label={APP_LAUNCHER_COPY.SEARCH_LABEL}
                  dataTestId={APP_LAUNCHER_TEST_IDS.SEARCH}/></div>
              <span className="catalog-content__count">{appCatalog.length} available</span>
            </div>
            {pinnedApps.length > 0 && (
              <section className="catalog-section" data-testid={APP_LAUNCHER_TEST_IDS.PINNED}>
                <div className="catalog-section__heading"><div><Pin size={16}/><h3>{APP_LAUNCHER_COPY.PINNED_TITLE}</h3></div><span>{APP_LAUNCHER_COPY.PINNED_DESCRIPTION}</span></div>
                <div className="catalog-app-grid">{pinnedApps.map(app => renderAppLauncherApp({ app, isPinned: true, onOpenChange, onTogglePin }))}</div>
              </section>
            )}
            {otherApps.length > 0 && (
              <section className="catalog-section" data-testid={APP_LAUNCHER_TEST_IDS.OTHER}>
                <div className="catalog-section__heading"><div><Grid2X2 size={16}/><h3>{APP_LAUNCHER_COPY.OTHER_TITLE}</h3></div><span>{APP_LAUNCHER_COPY.OTHER_DESCRIPTION}</span></div>
                <div className="catalog-app-grid">{otherApps.map(app => renderAppLauncherApp({ app, isPinned: false, onOpenChange, onTogglePin }))}</div>
              </section>
            )}
            {matching.length === 0 && <p className="catalog-empty" data-testid={APP_LAUNCHER_TEST_IDS.EMPTY}>{APP_LAUNCHER_COPY.EMPTY}</p>}
          </div>
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  );
};
