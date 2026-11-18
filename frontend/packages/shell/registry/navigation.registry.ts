import {
  Cable,
  LayoutDashboard,
  ShieldCheck,
  UsersRound,
  ContactRound,
} from 'lucide-react';

export const workspaceNavigation = [
  { label: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { label: 'People', path: '/people', icon: ContactRound },
  { label: 'Accounts', path: '/accounts', icon: UsersRound },
  { label: 'Integrations', path: '/integrations', icon: Cable },
  { label: 'Security', path: '/settings/security', icon: ShieldCheck },
] as const;
