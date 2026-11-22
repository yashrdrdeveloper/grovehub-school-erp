import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import { ShellContainer } from '@Shell';
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

// Apply the desktop preference before React mounts, without overriding system mode on phones.
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
if (!element) {
  throw new Error('Missing #root application mount point');
}

createRoot(element).render(
  <StrictMode>
    <ShellContainer />
  </StrictMode>,
);
