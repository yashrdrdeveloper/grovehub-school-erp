import { Bell, ClipboardCheck, Clock3, GraduationCap } from '@Components/icons';

// Layout preview only. No fabricated counts, events, tasks or school records.
export const DASHBOARD_METRICS = [
  {
    name: __('Tasks for you'),
    description: __('Your work queue'),
    icon: ClipboardCheck,
    tone: 'pink',
  },
  { name: __('Due today'), description: __('Deadlines & follow-ups'), icon: Clock3, tone: 'blue' },
  { name: __('Overdue'), description: __('Needs your attention'), icon: Bell, tone: 'amber' },
  { name: __('Students'), description: __('School community'), icon: GraduationCap, tone: 'green' },
] as const;
