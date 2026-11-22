import { Box } from '@Components/components/Box';
import { ArrowUpRight, CheckCircle2, CircleDashed, Inbox } from '@Components/icons';

import { DASHBOARD_WORKSPACES } from '../../constants/dashboard.workspaces';
import type { DashboardMainComponent } from './types/dashboardMain';

/** Preview placeholders for work queues and operational domains. */
export const DashboardMain: DashboardMainComponent = () => (
  <Box as="div" className="school-home__main-column">
    <Box as="section" className="school-home__section" aria-labelledby="home-queue">
      <Box as="div" className="school-home__section-heading">
        <Box as="div"><Box as="span" className="school-home__eyebrow">FOCUS / 02</Box><Box as="h2" id="home-queue">What needs attention</Box></Box>
        <Box as="span" className="school-home__hint">Your personal workspace</Box>
      </Box>
      <Box as="div" className="school-home__focus">
        <Box as="div" className="school-home__focus-heading">
          <Box as="div" className="school-home__focus-icon"><Inbox size={19}/></Box>
          <Box as="div"><Box as="h3">Pending for you</Box><Box as="p">Tasks, approvals and follow-ups assigned to you</Box></Box>
          <Box as="span" className="school-home__focus-count">Not connected</Box>
        </Box>
        <Box as="div" className="school-home__focus-content">
          <Box as="div" className="school-home__empty-art" aria-hidden="true"><Box as="span"/><Box as="span"/><Box as="span"/><CheckCircle2 size={24}/></Box>
          <Box as="strong">A clear place to start your day</Box>
          <Box as="p">Once Tasks is connected, your priorities will appear here, organized by what needs your attention first.</Box>
        </Box>
        <Box as="div" className="school-home__focus-footer">
          <Box as="span"><CircleDashed size={13}/> Waiting for Tasks service</Box><Box as="span">Workspace preview</Box>
        </Box>
      </Box>
    </Box>
    <Box as="section" className="school-home__section" aria-labelledby="home-operations">
      <Box as="div" className="school-home__section-heading">
        <Box as="div"><Box as="span" className="school-home__eyebrow">OPERATIONS / 03</Box><Box as="h2" id="home-operations">Around your school</Box></Box>
        <Box as="span" className="school-home__hint">Three work areas</Box>
      </Box>
      <Box as="div" className="school-home__operations">
        {DASHBOARD_WORKSPACES.map(({ title, hint, icon: Icon, tone }) => (
          <Box as="div" className={'school-home__operation school-home__tone--' + tone} key={title}>
            <Box as="div" className="school-home__operation-icon"><Icon size={18}/></Box>
            <Box as="div"><Box as="h3">{title}</Box><Box as="p">{hint}</Box><Box as="small">Ready for future service data</Box></Box>
            <ArrowUpRight size={16} className="school-home__operation-arrow"/>
          </Box>
        ))}
      </Box>
    </Box>
  </Box>
);
