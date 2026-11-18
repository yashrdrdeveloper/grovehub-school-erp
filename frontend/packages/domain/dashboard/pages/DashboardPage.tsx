import { ArrowRight, ArrowUpRight, BellRing, CalendarDays, CheckCheck, ChevronRight, CircleHelp, ClipboardList, Clock3, CreditCard, FileCheck2, GraduationCap, LayoutGrid, MessageSquareText, School, ShieldCheck, Sparkles, UsersRound } from 'lucide-react';
import { Link } from 'react-router-dom';

import { Card } from '@Components';

const stats = [
  { label: 'My tasks', description: 'Assignments & follow-ups', icon: ClipboardList, tone: 'violet' },
  { label: 'Due today', description: 'Work needing attention', icon: Clock3, tone: 'amber' },
  { label: 'Overdue', description: 'Past-due activities', icon: BellRing, tone: 'rose' },
  { label: 'Students', description: 'Enrolled student records', icon: GraduationCap, tone: 'blue' },
] as const;

const queues = [
  { title: 'Pending for you', description: 'Your assigned actions', icon: ClipboardList, tone: 'violet', empty: 'Your assigned tasks will appear here.', hint: 'Connect Tasks to see your queue.' },
  { title: 'Admissions & dues', description: 'Applications, fees and follow-ups', icon: CreditCard, tone: 'amber', empty: 'No admissions or fee data connected.', hint: 'Admissions and Finance are planned.' },
  { title: 'Today & upcoming', description: 'Events, holidays and schedules', icon: CalendarDays, tone: 'blue', empty: 'Upcoming school events will appear here.', hint: 'Calendar integration is planned.' },
  { title: 'Reminders', description: 'Things you need to remember', icon: BellRing, tone: 'rose', empty: 'Personal reminders will appear here.', hint: 'Reminder service is planned.' },
  { title: 'Complaints & requests', description: 'Open conversations and cases', icon: MessageSquareText, tone: 'cyan', empty: 'Requests will appear here when connected.', hint: 'Case management is planned.' },
  { title: 'Audit & compliance', description: 'Reviews, approvals and checks', icon: FileCheck2, tone: 'emerald', empty: 'Compliance reviews will appear here.', hint: 'Audit service is planned.' },
] as const;

const shortcuts = [
  { label: 'People directory', description: 'Find school members', href: '/people', icon: UsersRound },
  { label: 'Manage accounts', description: 'Permissions and access', href: '/accounts', icon: ShieldCheck },
  { label: 'Integrations', description: 'Connected applications', href: '/integrations', icon: LayoutGrid },
] as const;

export const DashboardPage = () => {
  const today = new Intl.DateTimeFormat(undefined, { weekday: 'long', month: 'long', day: 'numeric' }).format(new Date());

  return (
    <div className="page school-dashboard">
      <div className="dashboard-welcome">
        <div className="dashboard-welcome__copy">
          <span className="dashboard-welcome__date"><School size={15} /> RDR School <span aria-hidden="true">·</span> {today}</span>
          <h1>Good to see you <span aria-hidden="true">✳</span></h1>
          <p>Here is your school workspace. The most important things, all in one place.</p>
        </div>
        <Link className="dashboard-welcome__action" to="/accounts">Manage workspace <ArrowUpRight size={16}/></Link>
      </div>

      <div className="dashboard-banner">
        <div className="dashboard-banner__symbol"><Sparkles size={22}/></div>
        <div className="dashboard-banner__copy">
          <strong>Your school, connected</strong>
          <span>Start with people, accounts, and integrations. The operational dashboard fills in as you enable more school modules.</span>
        </div>
        <Link to="/integrations">Explore integrations <ArrowRight size={16}/></Link>
      </div>

      <section className="dashboard-block" aria-label="School summary">
        <div className="dashboard-block__heading">
          <div><span className="dashboard-kicker">AT A GLANCE</span><h2>Today's overview</h2></div>
          <span className="dashboard-quiet-label"><CircleHelp size={14}/> Live metrics will appear when services are connected</span>
        </div>
        <div className="metric-grid">
          {stats.map(({ label, description, icon: Icon, tone }) => (
            <Card key={label} className={`metric-card dashboard-tone dashboard-tone--${tone}`}>
              <div className="metric-card__top"><span>{label}</span><span className="metric-card__icon"><Icon size={19} strokeWidth={1.9}/></span></div>
              <strong aria-label={`${label}: not available`}>—</strong>
              <small>{description}</small>
            </Card>
          ))}
        </div>
      </section>

      <section className="dashboard-block" aria-labelledby="dashboard-work-title">
        <div className="dashboard-block__heading">
          <div><span className="dashboard-kicker">YOUR WORKSPACE</span><h2 id="dashboard-work-title">Everything that needs attention</h2></div>
          <span className="dashboard-quiet-label"><CheckCheck size={14}/> Ready for service data</span>
        </div>
        <div className="dashboard-panels">
          {queues.map(({ title, description, icon: Icon, tone, empty, hint }) => (
            <Card key={title} className={`dashboard-panel dashboard-tone dashboard-tone--${tone}`}>
              <div className="dashboard-panel__heading">
                <div className="dashboard-panel__icon"><Icon size={18}/></div>
                <div><h3>{title}</h3><small>{description}</small></div>
                <ChevronRight className="dashboard-panel__chevron" size={16}/>
              </div>
              <div className="dashboard-panel__body">
                <div className="dashboard-panel__empty"><div className="dashboard-panel__empty-symbol"><Icon size={21}/></div><strong>Nothing to show yet</strong><p>{empty}</p></div>
              </div>
              <div className="dashboard-panel__footer">{hint}</div>
            </Card>
          ))}
        </div>
      </section>

      <section className="dashboard-block" aria-labelledby="dashboard-links-title">
        <div className="dashboard-block__heading"><div><span className="dashboard-kicker">QUICK ACCESS</span><h2 id="dashboard-links-title">Jump back in</h2></div></div>
        <div className="dashboard-links">
          {shortcuts.map(({ label, description, href, icon: Icon }) => (
            <Link key={href} to={href} className="dashboard-quick-link">
              <span className="dashboard-quick-link__icon"><Icon size={19}/></span>
              <span><strong>{label}</strong><small>{description}</small></span>
              <ArrowUpRight size={17}/>
            </Link>
          ))}
        </div>
      </section>
    </div>
  );
};
