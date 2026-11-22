import { ArrowUpRight, CheckCircle2, CircleDashed, Inbox } from '@Components/icons';

import { DASHBOARD_WORKSPACES } from '../../constants/dashboard.workspaces';
import type { DashboardMainComponent } from './types/dashboardMain';

/** Preview placeholders for work queues and operational domains. */
export const DashboardMain: DashboardMainComponent = () => (
  <div className="school-home__main-column">
    <section className="school-home__section" aria-labelledby="home-queue">
      <div className="school-home__section-heading">
        <div><span className="school-home__eyebrow">FOCUS / 02</span><h2 id="home-queue">What needs attention</h2></div>
        <span className="school-home__hint">Your personal workspace</span>
      </div>
      <div className="school-home__focus">
        <div className="school-home__focus-heading">
          <div className="school-home__focus-icon"><Inbox size={19}/></div>
          <div><h3>Pending for you</h3><p>Tasks, approvals and follow-ups assigned to you</p></div>
          <span className="school-home__focus-count">Not connected</span>
        </div>
        <div className="school-home__focus-content">
          <div className="school-home__empty-art" aria-hidden="true"><span/><span/><span/><CheckCircle2 size={24}/></div>
          <strong>A clear place to start your day</strong>
          <p>Once Tasks is connected, your priorities will appear here, organized by what needs your attention first.</p>
        </div>
        <div className="school-home__focus-footer">
          <span><CircleDashed size={13}/> Waiting for Tasks service</span><span>Workspace preview</span>
        </div>
      </div>
    </section>
    <section className="school-home__section" aria-labelledby="home-operations">
      <div className="school-home__section-heading">
        <div><span className="school-home__eyebrow">OPERATIONS / 03</span><h2 id="home-operations">Around your school</h2></div>
        <span className="school-home__hint">Three work areas</span>
      </div>
      <div className="school-home__operations">
        {DASHBOARD_WORKSPACES.map(({ title, hint, icon: Icon, tone }) => (
          <div className={'school-home__operation school-home__tone--' + tone} key={title}>
            <div className="school-home__operation-icon"><Icon size={18}/></div>
            <div><h3>{title}</h3><p>{hint}</p><small>Ready for future service data</small></div>
            <ArrowUpRight size={16} className="school-home__operation-arrow"/>
          </div>
        ))}
      </div>
    </section>
  </div>
);
