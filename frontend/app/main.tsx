import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import { ShellContainer } from '@Shell';
import { resolveWorkspaceTheme } from '@Components/hooks/theme';

import '@Styles/index.css';
import './styles.css';
import '../packages/components/styles/primitives.css';

// Apply the desktop preference before React mounts, without overriding system mode on phones.
const applyInitialTheme = (): void => {
  try {
    const dark = window.matchMedia('(prefers-color-scheme: dark)');
    const phone = window.matchMedia('(max-width: 720px)');
    const apply = () => {
      const saved = window.localStorage.getItem('groovehub.theme');
      let preference: 'light' | 'dark' | 'system' = 'system';
      if (saved === 'light' || saved === 'dark') preference = saved;
      const resolved = resolveWorkspaceTheme({
        preference, systemDark: dark.matches, isMobile: phone.matches,
      });
      document.documentElement.dataset.theme = resolved;
      document.documentElement.style.colorScheme = resolved;
    };
    apply();
    dark.addEventListener('change', apply);
    phone.addEventListener('change', apply);
  } catch {
    const dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    if (dark) document.documentElement.dataset.theme = 'dark';
    else document.documentElement.dataset.theme = 'light';
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
