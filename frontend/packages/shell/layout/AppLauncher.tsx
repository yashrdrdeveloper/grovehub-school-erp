import { useState } from 'react';
import { BookOpen, BriefcaseBusiness, Building2, Cable, CalendarDays, ClipboardCheck, CreditCard, GraduationCap, Grid2X2, Megaphone, Search, Settings2, ShieldCheck, UsersRound, Workflow } from 'lucide-react';
import { Link } from 'react-router-dom';

import { Dialog } from '@Components';

import { workspaceNavigation } from '../registry';

type App = { name: string; description: string; path?: string; icon: typeof UsersRound; category: string };
const catalog: App[] = [
  { name: 'Dashboard', description: 'School overview', path: '/dashboard', icon: Building2, category: 'Workspace' },
  { name: 'People', description: 'Foundation directory', path: '/people', icon: UsersRound, category: 'Workspace' },
  { name: 'Accounts', description: 'User access', path: '/accounts', icon: ShieldCheck, category: 'Administration' },
  { name: 'Integrations', description: 'Connected services', path: '/integrations', icon: Cable, category: 'Administration' },
  { name: 'Security', description: 'Your sessions', path: '/settings/security', icon: Settings2, category: 'Administration' },
  { name: 'Admissions', description: 'Student admissions', icon: GraduationCap, category: 'Student life' },
  { name: 'Student lifecycle', description: 'Student records', icon: BookOpen, category: 'Student life' },
  { name: 'Attendance', description: 'Daily attendance', icon: ClipboardCheck, category: 'Student life' },
  { name: 'Fees & finance', description: 'Collections and dues', icon: CreditCard, category: 'Operations' },
  { name: 'Tasks', description: 'Work and follow-ups', icon: Workflow, category: 'Operations' },
  { name: 'Calendar', description: 'Academic schedule', icon: CalendarDays, category: 'Operations' },
  { name: 'Communication', description: 'Unified inbox', icon: Megaphone, category: 'Operations' },
  { name: 'HR', description: 'Staff and leave', icon: BriefcaseBusiness, category: 'Operations' },
];
const categories = ['All', 'Workspace', 'Student life', 'Operations', 'Administration'];

export interface AppLauncherProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}
export const AppLauncher = ({ open: controlled, onOpenChange }: AppLauncherProps) => {
  const [internalOpen, setInternalOpen] = useState(false);
  const [filter, setFilter] = useState('');
  const [category, setCategory] = useState('All');
  const open = controlled ?? internalOpen;
  const setOpen = onOpenChange ?? setInternalOpen;
  const apps = catalog.filter(app => (category === 'All' || app.category === category) &&
    `${app.name} ${app.description}`.toLowerCase().includes(filter.toLowerCase()));
  return (
    <>
      <button type="button" className="icon-button app-switcher-button" aria-label="Open application catalog"
        onClick={() => setOpen(true)}><Grid2X2 size={19}/></button>
      <Dialog open={open} onOpenChange={setOpen} title="App catalog" description="Your school applications, in one place.">
        <div className="catalog-search"><Search size={17}/><input value={filter} onChange={e => setFilter(e.target.value)} placeholder="Find an app..." aria-label="Search apps"/></div>
        <div className="catalog-tabs" aria-label="App categories">{categories.map(item => (
          <button key={item} type="button" className={category === item ? 'catalog-tab catalog-tab--active' : 'catalog-tab'}
            aria-pressed={category === item} onClick={() => setCategory(item)}>{item}</button>
        ))}</div>
        <div className="app-launcher-grid">
          {apps.map(({ name, description, path, icon: Icon }) => (
            path ? <Link key={name} to={path} className="app-launcher-item" onClick={() => setOpen(false)}>
              <span><Icon size={23} strokeWidth={1.7}/></span><strong>{name}</strong><small>{description}</small>
            </Link> : <div className="app-launcher-item app-launcher-item--disabled" key={name} aria-label={`${name}, planned`}>
              <span><Icon size={23} strokeWidth={1.7}/></span><strong>{name}</strong><small>Planned · {description}</small>
            </div>
          ))}
        </div>
        {apps.length === 0 && <p className="muted">No apps match your search.</p>}
        <p className="catalog-footnote">{workspaceNavigation.length} apps available · Other modules are planned and cannot be opened yet.</p>
      </Dialog>
    </>
  );
};
