import {
  BellRing, CalendarDays, CheckCheck, ChevronRight, CircleHelp,
  ClipboardList, Clock3, CreditCard, FileCheck2, GraduationCap,
  MessageSquareText, School, Sparkles,
} from 'lucide-react';

import { Card } from '@Components';

// Phase 0.2 visual placeholders only. Never represent these as live metrics.
const stats = [
  { label: 'My tasks', description: 'Assignments & follow-ups', icon: ClipboardList, tone: 'violet' },
  { label: 'Due today', description: 'Work needing attention', icon: Clock3, tone: 'amber' },
  { label: 'Overdue', description: 'Past-due activities', icon: BellRing, tone: 'rose' },
  { label: 'Students', description: 'Student records', icon: GraduationCap, tone: 'blue' },
] as const;

const queues = [
  { title: 'Pending for you', description: 'Your assigned work', icon: ClipboardList, tone: 'violet', empty: 'Assigned tasks will appear here.' },
  { title: 'Admissions & dues', description: 'Due today and overdue', icon: CreditCard, tone: 'amber', empty: 'Admission and fee follow-ups will appear here.' },
  { title: 'Today & upcoming', description: 'Holidays, events and time off', icon: CalendarDays, tone: 'blue', empty: 'School events and holidays will appear here.' },
  { title: 'Reminders', description: 'Your follow-ups', icon: BellRing, tone: 'rose', empty: 'Reminders will appear here.' },
  { title: 'Complaints & requests', description: 'Items needing review', icon: MessageSquareText, tone: 'cyan', empty: 'Requests and complaints will appear here.' },
  { title: 'Audit & compliance', description: 'Checks and reviews', icon: FileCheck2, tone: 'emerald', empty: 'Audits and compliance tasks will appear here.' },
] as const;

export const DashboardPage = () => {
  const today = new Intl.DateTimeFormat(undefined, {
    weekday: 'long', month: 'long', day: 'numeric',
  }).format(new Date());

  return (
    <div className="page school-dashboard">
      <header className="dashboard-welcome">
        <div className="dashboard-welcome__copy">
          <span className="dashboard-welcome__date"><School size={15}/> RDR School <span aria-hidden="true">·</span> {today}</span>
          <h1>Welcome to RDR School <span aria-hidden="true">✳</span></h1>
          <p>A single place for your school day. Your workspace will fill with live updates as modules are built.</p>
        </div>
      </header>

      <div className="dashboard-banner">
        <div className="dashboard-banner__symbol"><Sparkles size={22}/></div>
        <div className="dashboard-banner__copy">
          <strong>Your school at a glance</strong>
          <span>This is the dashboard layout preview. No student, admissions, tasks or financial data has been connected yet.</span>
        </div>
      </div>

      <section className="dashboard-block" aria-label="School summary">
        <div className="dashboard-block__heading">
          <div><span className="dashboard-kicker">AT A GLANCE</span><h2>Today's overview</h2></div>
          <span className="dashboard-quiet-label"><CircleHelp size={14}/> Preview · No live metrics</span>
        </div>
        <div className="metric-grid">
          {stats.map(({ label, description, icon: Icon, tone }) => (
            <Card key={label} className={`metric-card dashboard-tone dashboard-tone--${tone}`}>
              <div className="metric-card__top"><span>{label}</span><span className="metric-card__icon"><Icon size={19}/></span></div>
              <strong aria-label={`${label}: unavailable`}>—</strong>
              <small>{description}</small>
            </Card>
          ))}
        </div>
      </section>

      <section className="dashboard-block" aria-labelledby="dashboard-work-title">
        <div className="dashboard-block__heading">
          <div><span className="dashboard-kicker">YOUR WORKSPACE</span><h2 id="dashboard-work-title">Everything that needs attention</h2></div>
          <span className="dashboard-quiet-label"><CheckCheck size={14}/> Layout preview</span>
        </div>
        <div className="dashboard-panels">
          {queues.map(({ title, description, icon: Icon, tone, empty }) => (
            <Card key={title} className={`dashboard-panel dashboard-tone dashboard-tone--${tone}`}>
              <div className="dashboard-panel__heading">
                <div className="dashboard-panel__icon"><Icon size={18}/></div>
                <div><h3>{title}</h3><small>{description}</small></div>
                <ChevronRight className="dashboard-panel__chevron" size={16}/>
              </div>
              <div className="dashboard-panel__body">
                <div className="dashboard-panel__empty">
                  <div className="dashboard-panel__empty-symbol"><Icon size={21}/></div>
                  <strong>Coming soon</strong>
                  <p>{empty}</p>
                </div>
              </div>
              <div className="dashboard-panel__footer">Module not connected</div>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
};
