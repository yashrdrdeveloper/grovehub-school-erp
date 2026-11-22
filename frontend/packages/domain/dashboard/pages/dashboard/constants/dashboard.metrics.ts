import { Bell, ClipboardCheck, Clock3, GraduationCap } from '@Components/icons';

// Layout preview only. No fabricated counts, events, tasks or school records.
export const DASHBOARD_METRICS = [
  { name: 'Tasks for you', description: 'Your work queue', icon: ClipboardCheck, tone: 'pink' },
  { name: 'Due today', description: 'Deadlines & follow-ups', icon: Clock3, tone: 'blue' },
  { name: 'Overdue', description: 'Needs your attention', icon: Bell, tone: 'amber' },
  { name: 'Students', description: 'School community', icon: GraduationCap, tone: 'green' },
] as const;
