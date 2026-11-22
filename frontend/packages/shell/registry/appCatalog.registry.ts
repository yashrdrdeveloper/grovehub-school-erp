import {
  BookOpen, BriefcaseBusiness, Building2, Cable, CalendarDays, ClipboardCheck,
  CreditCard, GraduationCap, Megaphone, Settings2, ShieldCheck, UsersRound, Workflow,
} from 'lucide-react';

export type CatalogApp = {
  name: string;
  description: string;
  path?: string;
  icon: typeof UsersRound;
  category: string;
};

export const appCatalog: CatalogApp[] = [
  { name: 'Dashboard', description: 'Overview of your school', path: '/dashboard', icon: Building2, category: 'Workspace' },
  { name: 'People', description: 'School community directory', icon: UsersRound, category: 'Workspace' },
  { name: 'Accounts', description: 'Accounts and access', icon: ShieldCheck, category: 'Administration' },
  { name: 'Integrations', description: 'Connected services', icon: Cable, category: 'Administration' },
  { name: 'Security', description: 'Sessions and account security', icon: Settings2, category: 'Administration' },
  { name: 'Admissions', description: 'Applications and enrollment', icon: GraduationCap, category: 'Student life' },
  { name: 'Student lifecycle', description: 'Student journey and records', icon: BookOpen, category: 'Student life' },
  { name: 'Attendance', description: 'Student and staff attendance', icon: ClipboardCheck, category: 'Student life' },
  { name: 'Fees & finance', description: 'Collections and dues', icon: CreditCard, category: 'Operations' },
  { name: 'Tasks', description: 'Assignments and follow-ups', icon: Workflow, category: 'Operations' },
  { name: 'Calendar', description: 'Academic events and holidays', icon: CalendarDays, category: 'Operations' },
  { name: 'Communication', description: 'Calls, mail, SMS and WhatsApp', icon: Megaphone, category: 'Operations' },
  { name: 'HR', description: 'People operations and leave', icon: BriefcaseBusiness, category: 'Operations' },
];

export const defaultPinnedApps = ['/dashboard'];

