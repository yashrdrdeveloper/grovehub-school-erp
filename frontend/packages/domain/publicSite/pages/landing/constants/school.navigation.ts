import { SCHOOL_ROUTE } from './school.routes';

export const SCHOOL_NAV_ITEMS = [
  { label: 'Our school', to: SCHOOL_ROUTE.ABOUT },
  { label: 'Academics', to: SCHOOL_ROUTE.ACADEMICS },
  { label: 'Campus life', to: SCHOOL_ROUTE.CAMPUS },
  { label: 'Admissions', to: SCHOOL_ROUTE.ADMISSIONS },
  { label: 'Gallery', to: SCHOOL_ROUTE.GALLERY },
  { label: 'CBSE Disclosure', to: SCHOOL_ROUTE.DISCLOSURE },
] as const;
