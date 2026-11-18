import {
  ArrowUpRight, BellRing, CalendarDays, Check, ChevronRight,
  ClipboardList, Clock3, CreditCard, FileCheck2, GraduationCap,
  Layers3, MessageSquareText, School, Sparkles,
} from 'lucide-react';

const stats = [
  { label: 'My tasks', description: 'Assigned to you', icon: ClipboardList, tone: 'violet' },
  { label: 'Due today', description: 'Needs attention', icon: Clock3, tone: 'amber' },
  { label: 'Overdue', description: 'Past due', icon: BellRing, tone: 'rose' },
  { label: 'Students', description: 'School records', icon: GraduationCap, tone: 'blue' },
] as const;

const sections = [
  { title: 'Pending for you', description: 'Your assigned work', icon: ClipboardList, tone: 'violet', note: 'Tasks and assignments will appear here.' },
  { title: 'Admissions & dues', description: 'Due today and overdue', icon: CreditCard, tone: 'amber', note: 'Application and payment follow-ups will appear here.' },
  { title: 'Today & upcoming', description: 'Holidays, events and time off', icon: CalendarDays, tone: 'blue', note: 'Your school calendar will appear here.' },
  { title: 'Reminders', description: 'Your follow-ups', icon: BellRing, tone: 'rose', note: 'Personal reminders will appear here.' },
  { title: 'Complaints & requests', description: 'Items needing review', icon: MessageSquareText, tone: 'cyan', note: 'Requests and conversations will appear here.' },
  { title: 'Audit & compliance', description: 'Checks and reviews', icon: FileCheck2, tone: 'emerald', note: 'Reviews and compliance tasks will appear here.' },
] as const;

export const DashboardPage = () => {
  const date = new Intl.DateTimeFormat(undefined, { weekday: 'long', day: 'numeric', month: 'long' }).format(new Date());

  return (
    <div className="school-dashboard dashboard-v2">
      <header className="dashboard-v2__heading">
        <div>
          <div className="dashboard-v2__eyebrow"><School size={14}/> RDR SCHOOL <span className="dashboard-v2__eyebrow-divider"/> WORKSPACE</div>
          <h1>School overview<span className="dashboard-v2__heading-spark">✳</span></h1>
          <p>Everything important, organized for your school day.</p>
        </div>
        <div className="dashboard-v2__date"><CalendarDays size={16}/>{date}</div>
      </header>

      <section className="dashboard-v2__overview" aria-label="Overview">
        <div className="dashboard-v2__section-title">
          <div><span className="dashboard-v2__overline">01 / OVERVIEW</span><h2>At a glance</h2></div>
          <span className="dashboard-v2__preview"><span/> Layout preview</span>
        </div>
        <div className="dashboard-v2__metrics">
          {stats.map(({ label, description, icon: Icon, tone }) => (
            <div key={label} className={`dashboard-v2__metric dashboard-tone--${tone}`}>
              <div className="dashboard-v2__metric-top"><span>{label}</span><Icon size={17}/></div>
              <div className="dashboard-v2__metric-value" aria-label={`${label} not yet available`}>—</div>
              <div className="dashboard-v2__metric-bottom">{description}<ArrowUpRight size={13}/></div>
            </div>
          ))}
        </div>
      </section>

      <div className="dashboard-v2__content-grid">
        <section className="dashboard-v2__work" aria-label="School workspaces">
          <div className="dashboard-v2__section-title">
            <div><span className="dashboard-v2__overline">02 / WORKSPACE</span><h2>Work that matters</h2></div>
            <span className="dashboard-v2__section-count">06 sections</span>
          </div>
          <div className="dashboard-v2__panel-grid">
            {sections.map(({ title, description, icon: Icon, tone, note }, index) => (
              <article key={title} className={`dashboard-v2__panel dashboard-tone--${tone}`}>
                <div className="dashboard-v2__panel-head">
                  <span className="dashboard-v2__panel-icon"><Icon size={18} strokeWidth={1.85}/></span>
                  <span className="dashboard-v2__panel-index">{String(index + 1).padStart(2, '0')}</span>
                </div>
                <h3>{title}</h3>
                <p className="dashboard-v2__panel-description">{description}</p>
                <div className="dashboard-v2__panel-placeholder">
                  <span className="dashboard-v2__placeholder-line dashboard-v2__placeholder-line--wide"/>
                  <span className="dashboard-v2__placeholder-line"/>
                  <span className="dashboard-v2__placeholder-line dashboard-v2__placeholder-line--short"/>
                </div>
                <div className="dashboard-v2__panel-footer"><span>{note}</span><ChevronRight size={16}/></div>
              </article>
            ))}
          </div>
        </section>

        <aside className="dashboard-v2__aside" aria-label="Workspace status">
          <div className="dashboard-v2__section-title">
            <div><span className="dashboard-v2__overline">03 / STATUS</span><h2>Getting started</h2></div>
          </div>
          <div className="dashboard-v2__spotlight">
            <div className="dashboard-v2__spotlight-top"><div className="dashboard-v2__spotlight-icon"><Sparkles size={19}/></div><span>YOUR WORKSPACE</span></div>
            <h3>A better way to run your school.</h3>
            <p>The foundation is in place. This dashboard will become your daily command center as school modules are connected.</p>
            <div className="dashboard-v2__spotlight-bottom"><Check size={15}/> Layout ready for real data</div>
          </div>
          <div className="dashboard-v2__roadmap">
            <div className="dashboard-v2__roadmap-heading"><Layers3 size={17}/><strong>Workspace roadmap</strong></div>
            <div className="dashboard-v2__roadmap-item"><span className="dashboard-v2__roadmap-dot dashboard-v2__roadmap-dot--ready"/><div><strong>Identity & access</strong><small>Authentication foundation</small></div><span>01</span></div>
            <div className="dashboard-v2__roadmap-item"><span className="dashboard-v2__roadmap-dot"/><div><strong>School operations</strong><small>Tasks, schedules and records</small></div><span>02</span></div>
            <div className="dashboard-v2__roadmap-item"><span className="dashboard-v2__roadmap-dot"/><div><strong>Insights & reporting</strong><small>Connected operational data</small></div><span>03</span></div>
          </div>
        </aside>
      </div>
    </div>
  );
};
