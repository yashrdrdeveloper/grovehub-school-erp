import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import { ShellContainer } from '@Shell';

import './styles.css';

// Apply the saved preference before React mounts, including on login screens.
try {
  const preference = window.localStorage.getItem('groovehub.theme');
  const theme = preference === 'light' || preference === 'dark'
    ? preference
    : window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  document.documentElement.dataset.theme = theme;
  document.documentElement.style.colorScheme = theme;
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
    const current = window.localStorage.getItem('groovehub.theme');
    if (current !== 'light' && current !== 'dark') {
      document.documentElement.dataset.theme = event.matches ? 'dark' : 'light';
      document.documentElement.style.colorScheme = event.matches ? 'dark' : 'light';
    }
  });
} catch {
  document.documentElement.dataset.theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

const element = document.getElementById('root');
if (!element) {
  throw new Error('Missing #root application mount point');
}

createRoot(element).render(
  <StrictMode>
    <ShellContainer />
  </StrictMode>,
);
