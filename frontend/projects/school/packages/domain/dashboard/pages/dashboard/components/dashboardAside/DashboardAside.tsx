import { Button } from '@Components/components/Button';
import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './DashboardAside.module.css';
import sectionStyles from '../../DashboardSection.module.css';
import { Bell, CalendarDays, CircleDashed, Layers2 } from '@Components/icons';

import { getDashboardCalendarCopy } from './helpers/dashboardAside.getCalendarCopy';
import type { DashboardAsideComponent } from './types/dashboardAside';

/** Calendar selection and reminder placeholder, without invented event data. */
export const DashboardAside: DashboardAsideComponent = ({ date, calendar }) => {
  const today = calendar.view === 'today';
  const { calendarDay, calendarHeading, calendarDetail, calendarEmptyTitle } =
    getDashboardCalendarCopy(date, today);
  const todayClass = (today && styles['is-active']) || '';
  const upcomingClass = (!today && styles['is-active']) || '';

  return (
    <Box as={BOX_TAG.ASIDE} className={styles['school-home__side-column']}>
      <Box
        as={BOX_TAG.SECTION}
        className={sectionStyles['school-home__section']}
        aria-labelledby="home-calendar"
      >
        <Box as={BOX_TAG.DIV} className={sectionStyles['school-home__section-heading']}>
          <Box as={BOX_TAG.DIV}>
            <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__eyebrow']}>
              {__('CALENDAR / 04')}
            </Box>
            <Box as={BOX_TAG.H2} id="home-calendar">
              {__('School calendar')}
            </Box>
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles['school-home__calendar']}>
          <Box
            as={BOX_TAG.DIV}
            className={styles['school-home__calendar-tabs']}
            role="group"
            aria-label={__('Calendar view')}
          >
            <Button
              unstyled
              type="button"
              aria-pressed={today}
              className={todayClass}
              onClick={calendar.showToday}
            >
              {__('Today')}
            </Button>
            <Button
              unstyled
              type="button"
              aria-pressed={!today}
              className={upcomingClass}
              onClick={calendar.showUpcoming}
            >
              {__('Upcoming')}
            </Button>
          </Box>
          <Box as={BOX_TAG.DIV} className={styles['school-home__calendar-content']}>
            <Box as={BOX_TAG.DIV} className={styles['school-home__calendar-date']}>
              <Box as={BOX_TAG.STRONG}>{calendarDay}</Box>
              <Box as={BOX_TAG.SPAN}>
                {calendarHeading}
                <Box as={BOX_TAG.SMALL}>{calendarDetail}</Box>
              </Box>
            </Box>
            <Box as={BOX_TAG.DIV} className={styles['school-home__calendar-line']} />
            <Box as={BOX_TAG.DIV} className={styles['school-home__calendar-empty']}>
              <CalendarDays size={22} />
              <Box as={BOX_TAG.STRONG}>{calendarEmptyTitle}</Box>
              <Box as={BOX_TAG.SPAN}>
                {__(
                  'Holidays, school events and time off will appear here when Calendar is connected.',
                )}
              </Box>
            </Box>
          </Box>
          <Box as={BOX_TAG.DIV} className={styles['school-home__calendar-footer']}>
            <CircleDashed size={13} />
            {__('Calendar integration planned')}
          </Box>
        </Box>
      </Box>
      <Box
        as={BOX_TAG.SECTION}
        className={sectionStyles['school-home__section']}
        aria-labelledby="home-reminders"
      >
        <Box as={BOX_TAG.DIV} className={sectionStyles['school-home__section-heading']}>
          <Box as={BOX_TAG.DIV}>
            <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__eyebrow']}>
              {__('FOLLOW-UPS / 05')}
            </Box>
            <Box as={BOX_TAG.H2} id="home-reminders">
              {__('Reminders')}
            </Box>
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles['school-home__reminders']}>
          <Box as={BOX_TAG.DIV} className={styles['school-home__reminders-head']}>
            <Box as={BOX_TAG.SPAN} className={styles['school-home__reminders-icon']}>
              <Bell size={18} />
            </Box>
            <Box as={BOX_TAG.STRONG}>{__('Keep the important things close.')}</Box>
          </Box>
          <Box as={BOX_TAG.P}>
            {__(
              'Your reminders and personal follow-ups will live here—right beside your school calendar.',
            )}
          </Box>
          <Box as={BOX_TAG.DIV} className={styles['school-home__reminders-foot']}>
            <Layers2 size={14} />
            {__('Coming with Tasks')}
          </Box>
        </Box>
      </Box>
    </Box>
  );
};
