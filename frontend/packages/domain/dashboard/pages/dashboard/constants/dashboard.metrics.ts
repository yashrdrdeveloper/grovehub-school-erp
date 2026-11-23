import { Bell, ClipboardCheck, Clock3, GraduationCap } from '@Components/icons';

import { DASHBOARD_TONE } from './dashboard.tone';

// Layout preview only. No fabricated counts, events, tasks or school records.
export const DASHBOARD_METRICS = [
  {
    name: __('Tasks for you'),
    description: __('Your work queue'),
    icon: ClipboardCheck,
    tone: DASHBOARD_TONE.PINK,
  },
  {
    name: __('Due today'),
    description: __('Deadlines & follow-ups'),
    icon: Clock3,
    tone: DASHBOARD_TONE.BLUE,
  },
  {
    name: __('Overdue'),
    description: __('Needs your attention'),
    icon: Bell,
    tone: DASHBOARD_TONE.AMBER,
  },
  {
    name: __('Students'),
    description: __('School community'),
    icon: GraduationCap,
    tone: DASHBOARD_TONE.GREEN,
  },
] as const;
