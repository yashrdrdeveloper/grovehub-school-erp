import { _map } from '@Base/utilities/map';
import { cn } from '@Base/utilities/cn';
import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './DashboardMain.module.css';
import sectionStyles from '../../DashboardSection.module.css';
import { ArrowUpRight, CheckCircle2, CircleDashed, Inbox } from '@Components/icons';

import { DASHBOARD_WORKSPACES } from '../../constants/dashboard.workspaces';
import type { DashboardMainComponent } from './types/dashboardMain';

/** Preview placeholders for work queues and operational domains. */
export const DashboardMain: DashboardMainComponent = () => (
  <Box as={BOX_TAG.DIV} className={styles['school-home__main-column']}>
    <Box
      as={BOX_TAG.SECTION}
      className={sectionStyles['school-home__section']}
      aria-labelledby="home-queue"
    >
      <Box as={BOX_TAG.DIV} className={sectionStyles['school-home__section-heading']}>
        <Box as={BOX_TAG.DIV}>
          <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__eyebrow']}>
            FOCUS / 02
          </Box>
          <Box as={BOX_TAG.H2} id="home-queue">
            What needs attention
          </Box>
        </Box>
        <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__hint']}>
          Your personal workspace
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className={styles['school-home__focus']}>
        <Box as={BOX_TAG.DIV} className={styles['school-home__focus-heading']}>
          <Box as={BOX_TAG.DIV} className={styles['school-home__focus-icon']}>
            <Inbox size={19} />
          </Box>
          <Box as={BOX_TAG.DIV}>
            <Box as={BOX_TAG.H3}>Pending for you</Box>
            <Box as={BOX_TAG.P}>Tasks, approvals and follow-ups assigned to you</Box>
          </Box>
          <Box as={BOX_TAG.SPAN} className={styles['school-home__focus-count']}>
            Not connected
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles['school-home__focus-content']}>
          <Box as={BOX_TAG.DIV} className={styles['school-home__empty-art']} aria-hidden="true">
            <Box as={BOX_TAG.SPAN} />
            <Box as={BOX_TAG.SPAN} />
            <Box as={BOX_TAG.SPAN} />
            <CheckCircle2 size={24} />
          </Box>
          <Box as={BOX_TAG.STRONG}>A clear place to start your day</Box>
          <Box as={BOX_TAG.P}>
            Once Tasks is connected, your priorities will appear here, organized by what needs your
            attention first.
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles['school-home__focus-footer']}>
          <Box as={BOX_TAG.SPAN}>
            <CircleDashed size={13} /> Waiting for Tasks service
          </Box>
          <Box as={BOX_TAG.SPAN}>Workspace preview</Box>
        </Box>
      </Box>
    </Box>
    <Box
      as={BOX_TAG.SECTION}
      className={sectionStyles['school-home__section']}
      aria-labelledby="home-operations"
    >
      <Box as={BOX_TAG.DIV} className={sectionStyles['school-home__section-heading']}>
        <Box as={BOX_TAG.DIV}>
          <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__eyebrow']}>
            OPERATIONS / 03
          </Box>
          <Box as={BOX_TAG.H2} id="home-operations">
            Around your school
          </Box>
        </Box>
        <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__hint']}>
          Three work areas
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className={styles['school-home__operations']}>
        {_map(DASHBOARD_WORKSPACES, ({ title, hint, icon: Icon, tone }) => (
          <Box
            as={BOX_TAG.DIV}
            className={cn(
              styles['school-home__operation'],
              sectionStyles['school-home__tone--' + tone],
            )}
            key={title}
          >
            <Box as={BOX_TAG.DIV} className={styles['school-home__operation-icon']}>
              <Icon size={18} />
            </Box>
            <Box as={BOX_TAG.DIV}>
              <Box as={BOX_TAG.H3}>{title}</Box>
              <Box as={BOX_TAG.P}>{hint}</Box>
              <Box as={BOX_TAG.SMALL}>Ready for future service data</Box>
            </Box>
            <ArrowUpRight size={16} className={styles['school-home__operation-arrow']} />
          </Box>
        ))}
      </Box>
    </Box>
  </Box>
);
