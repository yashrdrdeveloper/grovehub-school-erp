import { Activity, ArrowUpRight, CalendarDays, CircleAlert, ClipboardList, Clock3, FileCheck2, GraduationCap, MessageSquareText, School, UsersRound, Wallet } from 'lucide-react';
import { Link } from 'react-router-dom';

import { Card, PageHeading } from '@Components';

const metrics = [
  { label: 'Pending tasks', icon: ClipboardList, detail: 'Task service not connected' },
  { label: 'Due today', icon: Clock3, detail: 'Admissions and fee reminders' },
  { label: 'Overdue', icon: CircleAlert, detail: 'No task data connected' },
  { label: 'Students', icon: GraduationCap, detail: 'Student lifecycle pending' },
];

const panels = [
  { title: 'Pending for you', subtitle: 'Your assigned work', icon: ClipboardList, message: 'Tasks will appear here when connected.', group: 'Work queue' },
  { title: 'Admissions & dues', subtitle: 'Due today and overdue', icon: Wallet, message: 'Admissions and finance data not connected.', group: 'Operations' },
  { title: 'Today & upcoming', subtitle: 'School schedule', icon: CalendarDays, message: 'Calendar events will appear here.', group: 'Schedule' },
  { title: 'Reminders', subtitle: 'Your follow-ups', icon: Clock3, message: 'No reminder service connected.', group: 'Schedule' },
  { title: 'Complaints & requests', subtitle: 'Cases requiring attention', icon: MessageSquareText, message: 'Requests will appear when the service is ready.', group: 'Operations' },
  { title: 'Audit & compliance', subtitle: 'Checks and pending reviews', icon: FileCheck2, message: 'Audit data is not connected yet.', group: 'Governance' },
];

export const DashboardPage = () => (
  <div className="page school-dashboard">
    <PageHeading eyebrow="RDR School · Operations" title="School overview" description="Your school, at a glance. Only connected services display live data." />
    <div className="dashboard-hero">
      <div>
        <span className="eyebrow eyebrow--light"><School size={13}/> RDR SCHOOL</span>
        <h2>Welcome to your workspace</h2>
        <p>Priorities, people, and school operations in one place.</p>
      </div>
      <Link to="/accounts" className="hero-link">Manage accounts <ArrowUpRight size={16}/></Link>
    </div>
    <div className="metric-grid">
      {metrics.map(({ label, icon: Icon, detail }) => (
        <Card key={label} className="metric-card">
          <div className="metric-card__top"><span>{label}</span><Icon size={18}/></div>
          <strong aria-label={`${label} unavailable`}>—</strong>
          <small>{detail}</small>
        </Card>
      ))}
    </div>
    <div className="dashboard-section-heading"><div><h2>Operations at a glance</h2><p>Work queues and schedules, organized by priority</p></div><span className="dashboard-data-indicator"><Activity size={13}/> Awaiting service data</span></div>
    <div className="dashboard-panels">
      {panels.map(({ title, subtitle, icon: Icon, message, group }) => (
        <Card key={title} className="dashboard-panel">
          <div className="dashboard-panel__heading"><div className="dashboard-panel__icon"><Icon size={18}/></div>
            <div><h3>{title}</h3><small>{subtitle}</small></div>
            <span className="dashboard-panel__group">{group}</span>
          </div>
          <div className="dashboard-panel__body"><div className="dashboard-panel__empty"><Icon size={22}/><p>{message}</p></div></div>
        </Card>
      ))}
    </div>
    <Card className="dashboard-shortcuts">
      <div><h3>Workspace shortcuts</h3><p className="muted">Manage the services available today</p></div>
      <div className="shortcut-list"><Link to="/people"><UsersRound size={17}/> People <ArrowUpRight size={14}/></Link><Link to="/accounts"><GraduationCap size={17}/> Accounts <ArrowUpRight size={14}/></Link></div>
    </Card>
  </div>
);
