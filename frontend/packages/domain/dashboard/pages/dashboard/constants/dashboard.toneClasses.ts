import sectionStyles from '../DashboardSection.module.css';

import { DASHBOARD_TONE } from './dashboard.tone';
import type { DashboardTone } from './types/dashboard.tone';

/** Exhaustive CSS module selection without computed class keys. */
export const DASHBOARD_TONE_CLASSES: Readonly<Record<DashboardTone, string>> = Object.freeze({
  [DASHBOARD_TONE.PINK]: sectionStyles['school-home__tone--pink'],
  [DASHBOARD_TONE.BLUE]: sectionStyles['school-home__tone--blue'],
  [DASHBOARD_TONE.AMBER]: sectionStyles['school-home__tone--amber'],
  [DASHBOARD_TONE.GREEN]: sectionStyles['school-home__tone--green'],
});
