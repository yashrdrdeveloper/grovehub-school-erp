import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import { installShellI18n } from '@Shell/containers/shellContainer/helpers/shellContainer.installI18n';
import {
  readThemePreference,
  resolveWorkspaceTheme,
  ResolvedTheme,
} from '@Components/hooks/useWorkspaceTheme';
import {
  WORKSPACE_THEME_DARK_QUERY,
  WORKSPACE_THEME_PHONE_QUERY,
} from '@Components/hooks/useWorkspaceTheme/constants/workspaceTheme.storage';

import '@Styles/index.css';
import './styles.css';
import '../packages/components/styles/primitives.css';

// Install the shell-owned global BEFORE dynamically importing any feature COPY constants.
// A static import of @Shell here would evaluate eager __('literal') calls too early.
const applyInitialTheme = (): void => {
  try {
    const dark = window.matchMedia(WORKSPACE_THEME_DARK_QUERY);
    const phone = window.matchMedia(WORKSPACE_THEME_PHONE_QUERY);
    const apply = () => {
      const preference = readThemePreference();
      const resolved = resolveWorkspaceTheme({
        preference,
        systemDark: dark.matches,
        isMobile: phone.matches,
      });
      document.documentElement.dataset.theme = resolved;
      document.documentElement.style.colorScheme = resolved;
    };
    apply();
    dark.addEventListener('change', apply);
    phone.addEventListener('change', apply);
  } catch {
    const dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    document.documentElement.dataset.theme = (dark && ResolvedTheme.DARK) || ResolvedTheme.LIGHT;
  }
};

applyInitialTheme();

const element = document.getElementById('root');
if (!element) throw new Error('Missing #root application mount point');

const startApplication = async (): Promise<void> => {
  try {
    installShellI18n();
    const { ShellContainer } = await import('@Shell');
    createRoot(element).render(
      <StrictMode>
        <ShellContainer />
      </StrictMode>,
    );
  } catch {
    // No global __() here: the bootstrap itself may have failed.
    element.textContent = 'Unable to load RDR School. Please refresh the page.';
  }
};

void startApplication();
