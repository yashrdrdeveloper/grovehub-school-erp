import { useState } from 'react';
import {
  ArrowUpRight, Bell, CalendarDays, CheckCircle2, ChevronRight, CircleDashed,
  ClipboardCheck, Clock3, CreditCard, GraduationCap, Inbox, Layers2,
  MessageCircle, School, ShieldCheck, Sparkles,
} from 'lucide-react';

// Layout-only preview: no pretend counts, trends, events, people or services.
const metrics = [
  { name: 'Tasks for you', description: 'Your work queue', icon: ClipboardCheck, tone: 'pink' },
  { name: 'Due today', description: 'Deadlines & follow-ups', icon: Clock3, tone: 'blue' },
  { name: 'Overdue', description: 'Needs your attention', icon: Bell, tone: 'amber' },
  { name: 'Students', description: 'School community', icon: GraduationCap, tone: 'green' },
] as const;

const workspaces = [
  { title: 'Admissions & dues', hint: 'Applications, payments and follow-ups', icon: CreditCard, tone: 'blue' },
  { title: 'Complaints & requests', hint: 'Conversations that need a response', icon: MessageCircle, tone: 'pink' },
  { title: 'Audit & compliance', hint: 'Checks, reviews and approvals', icon: ShieldCheck, tone: 'green' },
] as const;

export const DashboardPage = () => {
  const [calendarView, setCalendarView] = useState<'today' | 'upcoming'>('today');
  const date = new Date();
  const dayName = new Intl.DateTimeFormat(undefined, { weekday: 'long' }).format(date);
  const dateLabel = new Intl.DateTimeFormat(undefined, { day: 'numeric', month: 'long', year: 'numeric' }).format(date);

  return (
    <div className="school-home">
      <header className="school-home__header">
        <div>
          <div className="school-home__breadcrumb"><School size={14}/> RDR SCHOOL <ChevronRight size={13}/> OVERVIEW</div>
          <h1>Good {date.getHours() < 12 ? 'morning' : date.getHours() < 17 ? 'afternoon' : 'evening'}<span className="school-home__spark">✳</span></h1>
          <p>Your school day, in one clear view.</p>
        </div>
        <div className="school-home__date"><CalendarDays size={17}/><span><strong>{dayName}</strong><small>{dateLabel}</small></span></div>
      </header>

      <div className="school-home__notice" role="status">
        <div className="school-home__notice-mark"><Sparkles size={17}/></div>
        <div><strong>Your workspace is taking shape</strong><span>This is a preview of the dashboard layout. Your school's activity will appear as services are connected.</span></div>
        <span className="school-home__notice-state"><CircleDashed size={14}/> PREVIEW</span>
      </div>

      <section className="school-home__section" aria-labelledby="home-priorities">
        <div className="school-home__section-heading">
          <div><span className="school-home__eyebrow">YOUR DAY / 01</span><h2 id="home-priorities">At a glance</h2></div>
          <span className="school-home__hint">No live data yet</span>
        </div>
        <div className="school-home__stats">
          {metrics.map(({name,description,icon: Icon,tone}) => (
            <div className={`school-home__stat school-home__tone--${tone}`} key={name}>
              <div className="school-home__stat-label"><span>{name}</span><Icon size={18}/></div>
              <strong aria-label={`${name} not available`}>—</strong>
              <small>{description}</small>
              <div className="school-home__stat-track"><span/></div>
            </div>
          ))}
        </div>
      </section>

      <div className="school-home__layout">
        <div className="school-home__main-column">
          <section className="school-home__section" aria-labelledby="home-queue">
            <div className="school-home__section-heading">
              <div><span className="school-home__eyebrow">FOCUS / 02</span><h2 id="home-queue">What needs attention</h2></div>
              <span className="school-home__hint">Your personal workspace</span>
            </div>
            <div className="school-home__focus">
              <div className="school-home__focus-heading"><div className="school-home__focus-icon"><Inbox size={19}/></div><div><h3>Pending for you</h3><p>Tasks, approvals and follow-ups assigned to you</p></div><span className="school-home__focus-count">Not connected</span></div>
              <div className="school-home__focus-content">
                <div className="school-home__empty-art" aria-hidden="true"><span/><span/><span/><CheckCircle2 size={24}/></div>
                <strong>A clear place to start your day</strong>
                <p>Once Tasks is connected, your priorities will appear here, organized by what needs your attention first.</p>
              </div>
              <div className="school-home__focus-footer"><span><CircleDashed size={13}/> Waiting for Tasks service</span><span>Workspace preview</span></div>
            </div>
          </section>
          <section className="school-home__section" aria-labelledby="home-operations">
            <div className="school-home__section-heading">
              <div><span className="school-home__eyebrow">OPERATIONS / 03</span><h2 id="home-operations">Around your school</h2></div>
              <span className="school-home__hint">Three work areas</span>
            </div>
            <div className="school-home__operations">
              {workspaces.map(({title,hint,icon: Icon,tone}) => (
                <div className={`school-home__operation school-home__tone--${tone}`} key={title}>
                  <div className="school-home__operation-icon"><Icon size={18}/></div>
                  <div><h3>{title}</h3><p>{hint}</p><small>Ready for future service data</small></div>
                  <ArrowUpRight size={16} className="school-home__operation-arrow"/>
                </div>
              ))}
            </div>
          </section>
        </div>

        <aside className="school-home__side-column">
          <section className="school-home__section" aria-labelledby="home-calendar">
            <div className="school-home__section-heading">
              <div><span className="school-home__eyebrow">CALENDAR / 04</span><h2 id="home-calendar">School calendar</h2></div>
            </div>
            <div className="school-home__calendar">
              <div className="school-home__calendar-tabs" role="group" aria-label="Calendar view">
                <button type="button" aria-pressed={calendarView === 'today'} className={calendarView === 'today' ? 'is-active' : ''} onClick={()=>setCalendarView('today')}>Today</button>
                <button type="button" aria-pressed={calendarView === 'upcoming'} className={calendarView === 'upcoming' ? 'is-active' : ''} onClick={()=>setCalendarView('upcoming')}>Upcoming</button>
              </div>
              <div className="school-home__calendar-content">
                <div className="school-home__calendar-date"><strong>{calendarView === 'today' ? date.getDate() : '→'}</strong><span>{calendarView === 'today' ? dayName : 'Next on your calendar'}<small>{calendarView === 'today' ? dateLabel : 'Calendar preview'}</small></span></div>
                <div className="school-home__calendar-line"/>
                <div className="school-home__calendar-empty"><CalendarDays size={22}/><strong>{calendarView === 'today' ? 'Today’s schedule' : 'Upcoming events'}</strong><span>Holidays, school events and time off will appear here when Calendar is connected.</span></div>
              </div>
              <div className="school-home__calendar-footer"><CircleDashed size={13}/> Calendar integration planned</div>
            </div>
          </section>
          <section className="school-home__section" aria-labelledby="home-reminders">
            <div className="school-home__section-heading"><div><span className="school-home__eyebrow">FOLLOW-UPS / 05</span><h2 id="home-reminders">Reminders</h2></div></div>
            <div className="school-home__reminders">
              <div className="school-home__reminders-head"><span className="school-home__reminders-icon"><Bell size={18}/></span><strong>Keep the important things close.</strong></div>
              <p>Your reminders and personal follow-ups will live here—right beside your school calendar.</p>
              <div className="school-home__reminders-foot"><Layers2 size={14}/> Coming with Tasks</div>
            </div>
          </section>
        </aside>
      </div>
    </div>
  );
};
