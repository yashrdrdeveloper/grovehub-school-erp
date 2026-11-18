import { useState } from 'react';
import {
  ChevronDown, ChevronsLeft, ChevronsRight, Grid2X2, LayoutDashboard,
  Menu, Monitor, Moon, PinOff, Settings2, Sun, X,
} from 'lucide-react';
import { NavLink, Outlet } from 'react-router-dom';

import { AppLauncher, appCatalog, defaultPinnedApps } from './layout/AppLauncher';
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

const pinnedKey = 'groovehub.workspace.pinned.v1';
const collapsedKey = 'groovehub.workspace.sidebar.collapsed.v1';
const readPinned = (): string[] => {
  try {
    const raw = window.localStorage.getItem(pinnedKey);
    if (!raw) return [...defaultPinnedApps];
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [...defaultPinnedApps];
    const known = new Set(appCatalog.filter(app => app.path).map(app => app.path));
    return [...new Set(parsed.filter((path): path is string => typeof path === 'string' && known.has(path)))];
  } catch { return [...defaultPinnedApps]; }
};
const readCollapsed = (): boolean => {
  try { return window.localStorage.getItem(collapsedKey) === 'true'; } catch { return false; }
};

export const Shell = ({ accountLabel, schoolLabel, onLogout }: ShellProps) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(readCollapsed);
  const [quickSettings, setQuickSettings] = useState(false);
  const [launcherOpen, setLauncherOpen] = useState(false);
  const [pinned, setPinned] = useState<string[]>(readPinned);
  const { preference, setPreference } = useWorkspaceTheme();

  const toggleSidebar = () => {
    setSidebarCollapsed(previous => {
      const next = !previous;
      try { window.localStorage.setItem(collapsedKey, String(next)); } catch { /* local storage optional */ }
      return next;
    });
  };
  const togglePin = (path: string) => {
    setPinned(previous => {
      const next = previous.includes(path) ? previous.filter(item => item !== path) : [...previous, path];
      try { window.localStorage.setItem(pinnedKey, JSON.stringify(next)); } catch { /* local storage optional */ }
      return next;
    });
  };
  const navigation = pinned.flatMap(path => {
    const app = appCatalog.find(item => item.path === path);
    if (!app || !app.path) return [];
    const Icon = app.icon;
    return [
      <div className="sidebar__pinned-row" key={path}>
        <NavLink to={path} title={app.name} aria-label={app.name} onClick={() => setSidebarOpen(false)}
          className={({ isActive }) => `nav-item${isActive ? ' nav-item--active' : ''}`}>
          <Icon size={19} strokeWidth={1.8}/><span className="sidebar__label">{app.name}</span>
        </NavLink>
        <button type="button" className="sidebar__unpin" title={`Unpin ${app.name}`} aria-label={`Unpin ${app.name}`} onClick={() => togglePin(path)}><PinOff size={14}/></button>
      </div>,
    ];
  });

  return (
    <div className={`shell${sidebarCollapsed ? ' shell--collapsed' : ''}`}>
      <aside className={`sidebar${sidebarOpen ? ' sidebar--open' : ''}`} aria-label="School workspace sidebar">
        <div className="sidebar__header">
          <button type="button" className="icon-button app-switcher-button" title="Open app catalog" aria-label="Open app catalog" onClick={() => { setLauncherOpen(true); setSidebarOpen(false); }}><Grid2X2 size={19}/></button>
          <NavLink to="/dashboard" className="sidebar__brand" onClick={() => setSidebarOpen(false)} aria-label="GrooveHub home">
            <span className="brand-mark brand-mark--small">G</span><span className="sidebar__brand-name">groovehub<span className="brand-period">.</span></span>
          </NavLink>
          <button type="button" className="icon-button sidebar__close" aria-label="Close navigation" onClick={() => setSidebarOpen(false)}><X size={20}/></button>
        </div>
        <div className="sidebar__school" title={schoolLabel}>
          <div className="school-emblem">R</div>
          <span className="sidebar__school-details"><strong>{schoolLabel}</strong><small>School operating system</small></span>
        </div>
        <div className="sidebar__section-label">PINNED APPS</div>
        <nav aria-label="Pinned applications" className="sidebar__nav">
          {navigation.length > 0 ? navigation : <p className="sidebar__empty">Pin apps from the catalog to see them here.</p>}
        </nav>
        <button type="button" className="sidebar__catalog" title="Browse all applications" onClick={() => { setLauncherOpen(true); setSidebarOpen(false); }}>
          <Grid2X2 size={18}/><span className="sidebar__label">Browse all apps</span>
        </button>
        <div className="sidebar__bottom">
          <button type="button" className="sidebar__collapse" onClick={toggleSidebar} title={sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'} aria-label={sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}>
            {sidebarCollapsed ? <ChevronsRight size={17}/> : <ChevronsLeft size={17}/>}
            <span className="sidebar__label">{sidebarCollapsed ? 'Expand' : 'Collapse sidebar'}</span>
          </button>
          <button type="button" className="sidebar__settings" onClick={() => setQuickSettings(!quickSettings)} aria-expanded={quickSettings} title="Account and preferences">
            <span className="profile-avatar">{accountLabel.charAt(0).toUpperCase()}</span>
            <span className="sidebar__profile-details"><strong>{accountLabel}</strong><small>Account & preferences</small></span><Settings2 size={17} className="sidebar__settings-icon"/>
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
              <NavLink to="/settings/security" onClick={() => { setQuickSettings(false); setSidebarOpen(false); }}>Account security</NavLink>
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
            <button type="button" className="icon-button topbar__quick" onClick={() => { setQuickSettings(true); setSidebarOpen(true); }} aria-label="Open account preferences"><span className="profile-avatar profile-avatar--small">{accountLabel.charAt(0).toUpperCase()}</span><ChevronDown size={14}/></button>
          </div>
        </header>
        <main className="shell__main"><Outlet /></main>
        <nav className="mobile-nav" aria-label="Mobile navigation">
          <NavLink to="/dashboard" className={({isActive}) => `nav-item${isActive ? ' nav-item--active' : ''}`}><LayoutDashboard size={19}/><span>Home</span></NavLink>
          <button type="button" className="nav-item" onClick={() => setLauncherOpen(true)}><Grid2X2 size={19}/><span>Apps</span></button>
          <button type="button" className="nav-item" onClick={() => setSidebarOpen(true)}><Menu size={19}/><span>Menu</span></button>
        </nav>
      </div>
      <AppLauncher open={launcherOpen} onOpenChange={setLauncherOpen} pinned={pinned} onTogglePin={togglePin} sidebarCollapsed={sidebarCollapsed && !sidebarOpen}/>
    </div>
  );
};
