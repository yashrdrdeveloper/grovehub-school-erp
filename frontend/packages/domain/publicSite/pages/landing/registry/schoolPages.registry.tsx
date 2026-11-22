import type { ReactNode } from 'react';
import { SCHOOL_ROUTE } from '../constants/school.routes';
import { HomePage } from '../components/homeContent';
import { AboutPage, StoryPage } from '../components/aboutContent';
import { AcademicsPage, CampusPage, GalleryPage, AdmissionsPage } from '../components/campusContent';
import { DisclosurePage, ResourcesPage, ContactPage } from '../components/resourceContent';

export const pageContent: Record<string,()=>ReactNode> = {
  [SCHOOL_ROUTE.HOME]:HomePage,
  [SCHOOL_ROUTE.ABOUT]:AboutPage,
  [SCHOOL_ROUTE.HISTORY]:()=> <StoryPage topic="history"/>,
  [SCHOOL_ROUTE.DIRECTOR]:()=> <StoryPage topic="director"/>,
  [SCHOOL_ROUTE.PRINCIPAL]:()=> <StoryPage topic="principal"/>,
  [SCHOOL_ROUTE.PHILOSOPHY]:()=> <StoryPage topic="philosophy"/>,
  [SCHOOL_ROUTE.ACADEMICS]:AcademicsPage,
  [SCHOOL_ROUTE.CAMPUS]:CampusPage,
  [SCHOOL_ROUTE.ADMISSIONS]:AdmissionsPage,
  [SCHOOL_ROUTE.GALLERY]:GalleryPage,
  [SCHOOL_ROUTE.DISCLOSURE]:DisclosurePage,
  [SCHOOL_ROUTE.RESOURCES]:ResourcesPage,
  [SCHOOL_ROUTE.CONTACT]:ContactPage,
};

export const PUBLIC_SITE_PATHS = Object.keys(pageContent);
