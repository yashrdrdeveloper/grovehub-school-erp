import { Box } from '@Components/components/Box';
import {
  DialogRoot, DialogPortal, DialogOverlay, DialogContent,
  DialogTitle, DialogDescription, DialogClose,
} from '@Components/components/Dialog/parts';
import { InputControl } from '@Components/components/InputControl';
import { Grid2X2, Pin, Search, X } from '@Components/icons';
import { appCatalog } from '../../constants/shell.catalog';
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
          <Box as="div" className="catalog-content__header" data-testid={APP_LAUNCHER_TEST_IDS.HEADER}>
            <Box as="div" className="catalog-content__title">
              <Box as="span" className="catalog-content__mark"><Grid2X2 size={23}/></Box>
              <Box as="div"><DialogTitle>{APP_LAUNCHER_COPY.TITLE}</DialogTitle>
                <DialogDescription id="catalog-description">{APP_LAUNCHER_COPY.DESCRIPTION}</DialogDescription></Box>
            </Box>
            <DialogClose className="icon-button catalog-content__close" aria-label="Close app catalog"><X size={20}/></DialogClose>
          </Box>
          <Box as="div" className="catalog-content__body" data-testid={APP_LAUNCHER_TEST_IDS.BODY}>
            <Box as="div" className="catalog-content__tools">
              <Box as="div" className="catalog-search"><Search size={18}/>
                <InputControl unstyled value={filter} onChange={onFilterChange}
                  placeholder={APP_LAUNCHER_COPY.SEARCH_PLACEHOLDER} aria-label={APP_LAUNCHER_COPY.SEARCH_LABEL}
                  dataTestId={APP_LAUNCHER_TEST_IDS.SEARCH}/></Box>
              <Box as="span" className="catalog-content__count">{appCatalog.length} available</Box>
            </Box>
            {pinnedApps.length > 0 && (
              <Box as="section" className="catalog-section" data-testid={APP_LAUNCHER_TEST_IDS.PINNED}>
                <Box as="div" className="catalog-section__heading"><Box as="div"><Pin size={16}/><Box as="h3">{APP_LAUNCHER_COPY.PINNED_TITLE}</Box></Box><Box as="span">{APP_LAUNCHER_COPY.PINNED_DESCRIPTION}</Box></Box>
                <Box as="div" className="catalog-app-grid">{pinnedApps.map(app => renderAppLauncherApp({ app, isPinned: true, onOpenChange, onTogglePin }))}</Box>
              </Box>
            )}
            {otherApps.length > 0 && (
              <Box as="section" className="catalog-section" data-testid={APP_LAUNCHER_TEST_IDS.OTHER}>
                <Box as="div" className="catalog-section__heading"><Box as="div"><Grid2X2 size={16}/><Box as="h3">{APP_LAUNCHER_COPY.OTHER_TITLE}</Box></Box><Box as="span">{APP_LAUNCHER_COPY.OTHER_DESCRIPTION}</Box></Box>
                <Box as="div" className="catalog-app-grid">{otherApps.map(app => renderAppLauncherApp({ app, isPinned: false, onOpenChange, onTogglePin }))}</Box>
              </Box>
            )}
            {matching.length === 0 && <Box as="p" className="catalog-empty" data-testid={APP_LAUNCHER_TEST_IDS.EMPTY}>{APP_LAUNCHER_COPY.EMPTY}</Box>}
          </Box>
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  );
};
