import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './DashboardMain.module.css';
import sectionStyles from '../../DashboardSection.module.css';
import { CheckCircle2, CircleDashed, Inbox } from '@Components/icons';
import { renderDashboardMainWorkspace } from './renderer/dashboardMain.workspace';

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
            {__('FOCUS / 02')}
          </Box>
          <Box as={BOX_TAG.H2} id="home-queue">
            {__('What needs attention')}
          </Box>
        </Box>
        <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__hint']}>
          {__('Your personal workspace')}
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className={styles['school-home__focus']}>
        <Box as={BOX_TAG.DIV} className={styles['school-home__focus-heading']}>
          <Box as={BOX_TAG.DIV} className={styles['school-home__focus-icon']}>
            <Inbox size={19} />
          </Box>
          <Box as={BOX_TAG.DIV}>
            <Box as={BOX_TAG.H3}>{__('Pending for you')}</Box>
            <Box as={BOX_TAG.P}>{__('Tasks, approvals and follow-ups assigned to you')}</Box>
          </Box>
          <Box as={BOX_TAG.SPAN} className={styles['school-home__focus-count']}>
            {__('Not connected')}
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles['school-home__focus-content']}>
          <Box as={BOX_TAG.DIV} className={styles['school-home__empty-art']} aria-hidden="true">
            <Box as={BOX_TAG.SPAN} />
            <Box as={BOX_TAG.SPAN} />
            <Box as={BOX_TAG.SPAN} />
            <CheckCircle2 size={24} />
          </Box>
          <Box as={BOX_TAG.STRONG}>{__('A clear place to start your day')}</Box>
          <Box as={BOX_TAG.P}>
            {__(
              'Once Tasks is connected, your priorities will appear here, organized by what needs your attention first.',
            )}
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles['school-home__focus-footer']}>
          <Box as={BOX_TAG.SPAN}>
            <CircleDashed size={13} />
            {__('Waiting for Tasks service')}
          </Box>
          <Box as={BOX_TAG.SPAN}>{__('Workspace preview')}</Box>
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
            {__('OPERATIONS / 03')}
          </Box>
          <Box as={BOX_TAG.H2} id="home-operations">
            {__('Around your school')}
          </Box>
        </Box>
        <Box as={BOX_TAG.SPAN} className={sectionStyles['school-home__hint']}>
          {__('Three work areas')}
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className={styles['school-home__operations']}>
        {_map(DASHBOARD_WORKSPACES, renderDashboardMainWorkspace)}
      </Box>
    </Box>
  </Box>
);
