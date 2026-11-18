import { useState, type ReactNode } from 'react';
import { Bell, ChevronDown, Command, Menu, Moon, Search, Settings2, Sun, X } from 'lucide-react';
import { NavLink, Outlet } from 'react-router-dom';

import { Button } from '@Components';

import { workspaceNavigation } from './registry';
import { AppLauncher } from './layout/AppLauncher';

export interface ShellProps {
  accountLabel: string;
  schoolLabel: string;
  onLogout: () => void;
}

export const Shell = ({ accountLabel, schoolLabel, onLogout }: ShellProps) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [dark, setDark] = useState(false);
  const [quickSettings, setQuickSettings] = useState(false);

  const toggleTheme = () => {
    const next = !dark;
    setDark(next);
    document.documentElement.dataset.theme = next ? 'dark' : 'light';
  };

  const navigation: ReactNode = workspaceNavigation.map(({ label, path, icon: Icon }) => (
    <NavLink
      key={path}
      to={path}
      onClick={() => setSidebarOpen(false)}
      className={({ isActive }) => `nav-item${isActive ? ' nav-item--active' : ''}`}
    >
      <Icon size={19} strokeWidth={1.75} />
      <span>{label}</span>
    </NavLink>
  ));

  return (
    <div className="shell">
      <aside className={`sidebar${sidebarOpen ? ' sidebar--open' : ''}`}>
        <div className="sidebar__header">
          <AppLauncher />
          <NavLink to="/dashboard" className="sidebar__brand">
            <span className="brand-mark brand-mark--small">G</span>
            <span>groovehub<span className="brand-period">.</span></span>
          </NavLink>
          <button className="icon-button sidebar__close" aria-label="Close navigation" onClick={() => setSidebarOpen(false)}><X size={20} /></button>
        </div>
        <div className="sidebar__school">
          <div className="school-emblem">✦</div>
          <span><strong>{schoolLabel}</strong><small>School workspace</small></span>
          <ChevronDown size={15} />
        </div>
        <div className="sidebar__section-label">WORKSPACE</div>
        <nav aria-label="Primary navigation" className="sidebar__nav">{navigation}</nav>
        <div className="sidebar__bottom">
          <p className="sidebar__tip"><Command size={16} /> Your workspace, simplified.</p>
          <button className="sidebar__settings" onClick={() => setQuickSettings(!quickSettings)} aria-expanded={quickSettings}>
            <span className="profile-avatar">{accountLabel.charAt(0).toUpperCase()}</span>
            <span><strong>Quick settings</strong><small>Preferences & account</small></span>
            <Settings2 size={17} />
          </button>
          {quickSettings ? (
            <div className="settings-popover">
              <strong>{accountLabel}</strong>
              <NavLink to="/settings/security" onClick={() => setQuickSettings(false)}>Account security</NavLink>
              <button onClick={toggleTheme}>{dark ? 'Light theme' : 'Dark theme'}</button>
              <button onClick={onLogout}>Sign out</button>
            </div>
          ) : null}
        </div>
      </aside>
      {sidebarOpen ? <button className="sidebar-backdrop" aria-label="Dismiss navigation" onClick={() => setSidebarOpen(false)} /> : null}
      <div className="shell__content">
        <header className="topbar">
          <div className="topbar__left">
            <button className="icon-button topbar__menu" onClick={() => setSidebarOpen(true)} aria-label="Open menu"><Menu size={20} /></button>
            <span className="topbar__breadcrumb">{schoolLabel}</span>
            <span className="topbar__separator">/</span>
            <span className="topbar__muted">Workspace</span>
          </div>
          <div className="topbar__right">
            <button className="topbar__search" aria-label="Search (coming with global search)" disabled><Search size={17} /><span>Search</span><kbd>⌘ K</kbd></button>
            <button className="icon-button" aria-label="Toggle theme" onClick={toggleTheme}>{dark ? <Sun size={19} /> : <Moon size={19} />}</button>
            <button className="icon-button" aria-label="Notifications (coming with task service)" disabled><Bell size={19} /></button>
            <Button variant="ghost" className="topbar__quick" onClick={() => setQuickSettings(!quickSettings)}><span className="profile-avatar profile-avatar--small">{accountLabel.charAt(0).toUpperCase()}</span><ChevronDown size={14} /></Button>
          </div>
        </header>
        <main className="shell__main"><Outlet /></main>
        <nav className="mobile-nav" aria-label="Mobile navigation">{navigation}</nav>
      </div>
    </div>
  );
};
