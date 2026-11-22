import { Bell, CalendarDays, CircleDashed, Layers2 } from '@Components/icons';

import type { DashboardAsideComponent } from './types/dashboardAside';

/** Calendar selection and reminder placeholder, without invented event data. */
export const DashboardAside: DashboardAsideComponent = ({ date, calendar }) => {
  const today = calendar.view === 'today';
  const calendarDay = today ? String(date.dayOfMonth) : '→';
  const calendarHeading = today ? date.dayName : 'Next on your calendar';
  const calendarDetail = today ? date.dateLabel : 'Calendar preview';
  const calendarEmptyTitle = today ? 'Today’s schedule' : 'Upcoming events';

  return (
    <aside className="school-home__side-column">
      <section className="school-home__section" aria-labelledby="home-calendar">
        <div className="school-home__section-heading">
          <div><span className="school-home__eyebrow">CALENDAR / 04</span><h2 id="home-calendar">School calendar</h2></div>
        </div>
        <div className="school-home__calendar">
          <div className="school-home__calendar-tabs" role="group" aria-label="Calendar view">
            <button type="button" aria-pressed={today} className={today ? 'is-active' : ''}
              onClick={calendar.showToday}>Today</button>
            <button type="button" aria-pressed={!today} className={!today ? 'is-active' : ''}
              onClick={calendar.showUpcoming}>Upcoming</button>
          </div>
          <div className="school-home__calendar-content">
            <div className="school-home__calendar-date">
              <strong>{calendarDay}</strong>
              <span>{calendarHeading}<small>{calendarDetail}</small></span>
            </div>
            <div className="school-home__calendar-line"/>
            <div className="school-home__calendar-empty">
              <CalendarDays size={22}/>
              <strong>{calendarEmptyTitle}</strong>
              <span>Holidays, school events and time off will appear here when Calendar is connected.</span>
            </div>
          </div>
          <div className="school-home__calendar-footer"><CircleDashed size={13}/> Calendar integration planned</div>
        </div>
      </section>
      <section className="school-home__section" aria-labelledby="home-reminders">
        <div className="school-home__section-heading">
          <div><span className="school-home__eyebrow">FOLLOW-UPS / 05</span><h2 id="home-reminders">Reminders</h2></div>
        </div>
        <div className="school-home__reminders">
          <div className="school-home__reminders-head">
            <span className="school-home__reminders-icon"><Bell size={18}/></span>
            <strong>Keep the important things close.</strong>
          </div>
          <p>Your reminders and personal follow-ups will live here—right beside your school calendar.</p>
          <div className="school-home__reminders-foot"><Layers2 size={14}/> Coming with Tasks</div>
        </div>
      </section>
    </aside>
  );
};
