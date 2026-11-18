import { ArrowUpRight, CalendarDays, CircleCheck, ClipboardList, UsersRound } from 'lucide-react';
import { Link } from 'react-router-dom';

import { Card, PageHeading } from '@Components';

const metrics = [
  { label: 'Tasks', icon: ClipboardList, value: '—', detail: 'Connect Task service' },
  { label: 'Due today', icon: CircleCheck, value: '—', detail: 'Connect admissions & finance' },
  { label: 'Students', icon: UsersRound, value: '—', detail: 'Connect student lifecycle' },
  { label: 'Upcoming', icon: CalendarDays, value: '—', detail: 'Connect school calendar' },
];

export const DashboardPage = () => (
  <div className="page">
    <PageHeading eyebrow="Workspace overview" title="Dashboard" description="Everything your school needs, at a glance." />
    <div className="dashboard-hero">
      <div>
        <span className="eyebrow eyebrow--light">YOUR SCHOOL, ONE PLACE</span>
        <h2>Good to have you here.</h2>
        <p>Manage people, connections and your school operations from one workspace.</p>
      </div>
      <Link to="/accounts" className="hero-link">Manage accounts <ArrowUpRight size={16} /></Link>
    </div>
    <div className="metric-grid">
      {metrics.map(({ label, icon: Icon, value, detail }) => (
        <Card key={label} className="metric-card">
          <div className="metric-card__top"><span>{label}</span><Icon size={18} /></div>
          <strong>{value}</strong>
          <small>{detail}</small>
        </Card>
      ))}
    </div>
    <div className="dashboard-lower">
      <Card>
        <div className="section-heading"><h3>Pending for you</h3><span className="muted">Live data when the Task API is available</span></div>
        <div className="dashboard-placeholder"><ClipboardList size={23} /><p>No task data connected yet.</p></div>
      </Card>
      <Card>
        <div className="section-heading"><h3>Today & upcoming</h3><span className="muted">School calendar</span></div>
        <div className="dashboard-placeholder"><CalendarDays size={23} /><p>Calendar integration coming with the school service.</p></div>
      </Card>
    </div>
  </div>
);
