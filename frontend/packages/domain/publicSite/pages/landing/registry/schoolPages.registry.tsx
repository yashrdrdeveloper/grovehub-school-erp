import type { ReactNode } from 'react';
import { HomePage } from '../components/homeContent';
import { AboutPage, StoryPage } from '../components/aboutContent';
import { AcademicsPage, CampusPage, GalleryPage, AdmissionsPage } from '../components/campusContent';
import { DisclosurePage, ResourcesPage, ContactPage } from '../components/resourceContent';

export const pageContent: Record<string,()=>ReactNode> = {
  '/':HomePage,
  '/about':AboutPage,
  '/about/history':()=> <StoryPage topic="history"/>,
  '/about/director':()=> <StoryPage topic="director"/>,
  '/about/principal':()=> <StoryPage topic="principal"/>,
  '/about/philosophy':()=> <StoryPage topic="philosophy"/>,
  '/academics':AcademicsPage,
  '/campus':CampusPage,
  '/admissions':AdmissionsPage,
  '/gallery':GalleryPage,
  '/mandatory-public-disclosure':DisclosurePage,
  '/resources':ResourcesPage,
  '/contact':ContactPage,
};

export const PUBLIC_SITE_PATHS = Object.keys(pageContent);
