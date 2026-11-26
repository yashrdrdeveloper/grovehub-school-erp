import { useEffect, useRef } from 'react';

/** Focus the full-page phone menu and restore its trigger when the menu closes. */
export const useMobileSidebarFocus = (open: boolean) => {
  const sidebar = useRef<HTMLElement>(null);

  useEffect(() => {
    if (!open || typeof window.matchMedia !== 'function' || !window.matchMedia('(max-width: 720px)').matches) return;
    const element = sidebar.current;
    if (!element) return;

    element.querySelector<HTMLButtonElement>('.sidebar__close')?.focus();

    return () => {
      if (!element.isConnected) return;
      const active = document.activeElement;
      if (active && active !== document.body && !element.contains(active)) return;
      const trigger = document.getElementById('workspace-mobile-menu');
      if (trigger && window.getComputedStyle(trigger).visibility !== 'hidden') trigger.focus();
    };
  }, [open]);

  return sidebar;
};
