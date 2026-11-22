import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { Bell, CalendarDays, CircleDashed, Layers2 } from '@Components/icons';

import { getDashboardCalendarCopy } from './helpers/dashboardAside.getCalendarCopy';
import type { DashboardAsideComponent } from './types/dashboardAside';

/** Calendar selection and reminder placeholder, without invented event data. */
export const DashboardAside: DashboardAsideComponent = ({ date, calendar }) => {
  const today = calendar.view === 'today';
  const { calendarDay, calendarHeading, calendarDetail, calendarEmptyTitle } = getDashboardCalendarCopy(date, today);
  let todayClass = '';
  let upcomingClass = '';
  if (today) todayClass = 'is-active';
  else upcomingClass = 'is-active';

  return (
    <Box as="aside" className="school-home__side-column">
      <Box as="section" className="school-home__section" aria-labelledby="home-calendar">
        <Box as="div" className="school-home__section-heading">
          <Box as="div"><Box as="span" className="school-home__eyebrow">CALENDAR / 04</Box><Box as="h2" id="home-calendar">School calendar</Box></Box>
        </Box>
        <Box as="div" className="school-home__calendar">
          <Box as="div" className="school-home__calendar-tabs" role="group" aria-label="Calendar view">
            <Button unstyled type="button" aria-pressed={today} className={todayClass}
              onClick={calendar.showToday}>Today</Button>
            <Button unstyled type="button" aria-pressed={!today} className={upcomingClass}
              onClick={calendar.showUpcoming}>Upcoming</Button>
          </Box>
          <Box as="div" className="school-home__calendar-content">
            <Box as="div" className="school-home__calendar-date">
              <Box as="strong">{calendarDay}</Box>
              <Box as="span">{calendarHeading}<Box as="small">{calendarDetail}</Box></Box>
            </Box>
            <Box as="div" className="school-home__calendar-line"/>
            <Box as="div" className="school-home__calendar-empty">
              <CalendarDays size={22}/>
              <Box as="strong">{calendarEmptyTitle}</Box>
              <Box as="span">Holidays, school events and time off will appear here when Calendar is connected.</Box>
            </Box>
          </Box>
          <Box as="div" className="school-home__calendar-footer"><CircleDashed size={13}/> Calendar integration planned</Box>
        </Box>
      </Box>
      <Box as="section" className="school-home__section" aria-labelledby="home-reminders">
        <Box as="div" className="school-home__section-heading">
          <Box as="div"><Box as="span" className="school-home__eyebrow">FOLLOW-UPS / 05</Box><Box as="h2" id="home-reminders">Reminders</Box></Box>
        </Box>
        <Box as="div" className="school-home__reminders">
          <Box as="div" className="school-home__reminders-head">
            <Box as="span" className="school-home__reminders-icon"><Bell size={18}/></Box>
            <Box as="strong">Keep the important things close.</Box>
          </Box>
          <Box as="p">Your reminders and personal follow-ups will live here—right beside your school calendar.</Box>
          <Box as="div" className="school-home__reminders-foot"><Layers2 size={14}/> Coming with Tasks</Box>
        </Box>
      </Box>
    </Box>
  );
};
