import { useState } from 'react';
import { Bell, ChevronDown, Grid2X2, Menu, Monitor, Moon, Search, Settings2, Sun, X } from 'lucide-react';
import { NavLink, Outlet } from 'react-router-dom';

import { Button } from '@Components';

import { workspaceNavigation } from './registry';
import { AppLauncher } from './layout/AppLauncher';
import { useWorkspaceTheme, type ThemePreference } from './theme';

export interface ShellProps {
  accountLabel: string;
  schoolLabel: string;
  onLogout: () => void;
}

const choices: { value: ThemePreference; label: string; icon: typeof Sun }[] = [
  { value: 'system', label: 'System', icon: Monitor },
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
];

export const Shell = ({ accountLabel, schoolLabel, onLogout }: ShellProps) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [quickSettings, setQuickSettings] = useState(false);
  const [launcherOpen, setLauncherOpen] = useState(false);
  const { preference, setPreference } = useWorkspaceTheme();
  const navigation = workspaceNavigation.map(({ label, path, icon: Icon }) => (
    <NavLink key={path} to={path} onClick={() => setSidebarOpen(false)}
      className={({ isActive }) => `nav-item${isActive ? ' nav-item--active' : ''}`}>
      <Icon size={19} strokeWidth={1.8} /><span>{label}</span>
    </NavLink>
  ));
  return (
    <div className="shell">
      <aside className={`sidebar${sidebarOpen ? ' sidebar--open' : ''}`}>
        <div className="sidebar__header">
          <AppLauncher open={launcherOpen} onOpenChange={setLauncherOpen} />
          <NavLink to="/dashboard" className="sidebar__brand" onClick={() => setSidebarOpen(false)}>
            <span className="brand-mark brand-mark--small">G</span><span>groovehub<span className="brand-period">.</span></span>
          </NavLink>
          <button type="button" className="icon-button sidebar__close" aria-label="Close navigation" onClick={() => setSidebarOpen(false)}><X size={20}/></button>
        </div>
        <div className="sidebar__school">
          <div className="school-emblem">R</div>
          <span><strong>{schoolLabel}</strong><small>School operating system</small></span>
        </div>
        <div className="sidebar__section-label">PINNED WORKSPACE</div>
        <nav aria-label="Primary navigation" className="sidebar__nav">{navigation}</nav>
        <button className="sidebar__catalog" onClick={() => { setLauncherOpen(true); setSidebarOpen(false); }}>
          <Grid2X2 size={18}/> Browse app catalog
        </button>
        <div className="sidebar__bottom">
          <button type="button" className="sidebar__settings" onClick={() => setQuickSettings(!quickSettings)} aria-expanded={quickSettings}>
            <span className="profile-avatar">{accountLabel.charAt(0).toUpperCase()}</span>
            <span><strong>{accountLabel}</strong><small>Account & preferences</small></span><Settings2 size={17}/>
          </button>
          {quickSettings && (
            <div className="settings-popover">
              <strong>Appearance</strong>
              <div className="theme-choices" role="group" aria-label="Appearance">
                {choices.map(({ value, label, icon: Icon }) => (
                  <button key={value} type="button" aria-pressed={preference === value}
                    className={preference === value ? 'theme-choice theme-choice--active' : 'theme-choice'}
                    onClick={() => setPreference(value)}><Icon size={16}/>{label}</button>
                ))}
              </div>
              <NavLink to="/settings/security" onClick={() => setQuickSettings(false)}>Account security</NavLink>
              <button type="button" onClick={onLogout}>Sign out</button>
            </div>
          )}
        </div>
      </aside>
      {sidebarOpen && <button type="button" className="sidebar-backdrop" aria-label="Dismiss navigation" onClick={() => setSidebarOpen(false)}/>}
      <div className="shell__content">
        <header className="topbar">
          <div className="topbar__left">
            <button type="button" className="icon-button topbar__menu" onClick={() => setSidebarOpen(true)} aria-label="Open menu"><Menu size={20}/></button>
            <span className="topbar__breadcrumb">{schoolLabel}</span><span className="topbar__separator">/</span><span className="topbar__muted">Workspace</span>
          </div>
          <div className="topbar__right">
            <button type="button" className="topbar__search" aria-label="Global search not available yet" disabled><Search size={17}/><span>Search</span><kbd>⌘ K</kbd></button>
            <button type="button" className="icon-button" aria-label="Notifications not available yet" disabled><Bell size={19}/></button>
            <Button variant="ghost" className="topbar__quick" onClick={() => setQuickSettings(!quickSettings)}><span className="profile-avatar profile-avatar--small">{accountLabel.charAt(0).toUpperCase()}</span><ChevronDown size={14}/></Button>
          </div>
        </header>
        <main className="shell__main"><Outlet /></main>
        <nav className="mobile-nav" aria-label="Mobile navigation">
          <NavLink to="/dashboard" className={({isActive}) => `nav-item${isActive ? ' nav-item--active' : ''}`}><span className="mobile-nav__icon">⌂</span><span>Home</span></NavLink>
          <NavLink to="/people" className={({isActive}) => `nav-item${isActive ? ' nav-item--active' : ''}`}><span className="mobile-nav__icon">♙</span><span>People</span></NavLink>
          <button type="button" className="nav-item" onClick={() => setLauncherOpen(true)}><Grid2X2 size={19}/><span>Apps</span></button>
          <NavLink to="/accounts" className={({isActive}) => `nav-item${isActive ? ' nav-item--active' : ''}`}><span className="mobile-nav__icon">◎</span><span>Accounts</span></NavLink>
          <button type="button" className="nav-item" onClick={() => setSidebarOpen(true)}><Menu size={19}/><span>Menu</span></button>
        </nav>
      </div>
    </div>
  );
};
