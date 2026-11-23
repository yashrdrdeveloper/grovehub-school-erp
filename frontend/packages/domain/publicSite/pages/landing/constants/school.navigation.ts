import { SCHOOL_ROUTE } from './school.routes';

export const SCHOOL_NAV_ITEMS = [
  { label: __('Our school'), to: SCHOOL_ROUTE.ABOUT },
  { label: __('Academics'), to: SCHOOL_ROUTE.ACADEMICS },
  { label: __('Campus life'), to: SCHOOL_ROUTE.CAMPUS },
  { label: __('Admissions'), to: SCHOOL_ROUTE.ADMISSIONS },
  { label: __('Gallery'), to: SCHOOL_ROUTE.GALLERY },
  { label: __('CBSE Disclosure'), to: SCHOOL_ROUTE.DISCLOSURE },
] as const;
