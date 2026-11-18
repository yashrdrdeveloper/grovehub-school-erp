import { useEffect, useState } from 'react';

export type ThemePreference = 'system' | 'light' | 'dark';
const storageKey = 'groovehub.theme';

function storedPreference(): ThemePreference {
  try {
    const value = window.localStorage.getItem(storageKey);
    return value === 'light' || value === 'dark' ? value : 'system';
  } catch {
    return 'system';
  }
}

export function useWorkspaceTheme() {
  const [preference, setPreference] = useState<ThemePreference>(storedPreference);

  useEffect(() => {
    const media = window.matchMedia('(prefers-color-scheme: dark)');
    const apply = () => {
      const resolved = preference === 'system' ? (media.matches ? 'dark' : 'light') : preference;
      document.documentElement.dataset.theme = resolved;
      document.documentElement.style.colorScheme = resolved;
    };
    apply();
    media.addEventListener('change', apply);
    return () => media.removeEventListener('change', apply);
  }, [preference]);

  const update = (next: ThemePreference) => {
    try { window.localStorage.setItem(storageKey, next); } catch { /* storage may be unavailable */ }
    setPreference(next);
  };

  return { preference, setPreference: update };
}
