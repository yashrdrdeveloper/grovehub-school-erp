import type { ReactNode } from 'react';
import type { SchoolPanel } from '../components/publicSheet';
import { HomePage } from '../components/homeContent';
import { AboutPage, StoryPage } from '../components/aboutContent';
import { AcademicsPage, CampusPage, GalleryPage, AdmissionsPage } from '../components/campusContent';
import { DisclosurePage, ResourcesPage, ContactPage } from '../components/resourceContent';

const pageContent: Record<string,()=>ReactNode> = {
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

const pathPanel: Record<string,SchoolPanel> = {
  '/about':'about',
  '/about/history':'history',
  '/about/director':'director',
  '/about/principal':'principal',
  '/about/philosophy':'philosophy',
  '/academics':'academics',
  '/campus':'campus',
  '/admissions':'enquiry',
  '/gallery':'gallery',
  '/resources':'resources',
  '/contact':'contact',
};

const pageMetadata: Record<string,{title:string;description:string}> = {
  '/':{title:'RDR School | Rao Deena Ram Vidya Vihar, Haluhera',description:'Discover Rao Deena Ram Vidya Vihar School in Haluhera, Rewari. Learn about school life, academic programmes, facilities and admissions.'},
  '/about':{title:'About RDR School | History, Vision & Leadership',description:'Explore RDR School history since 2006, its philosophy, director and principal messages in Haluhera, Rewari.'},
  '/about/history':{title:'School History | RDR School Haluhera',description:'Discover how RDR School was founded in April 2006 in Haluhera, Rewari, by Shri Lal Singh Yadav.'},
  '/about/director':{title:"Director's Message | RDR School",description:"Discover Shri Lal Singh Yadav's published message about all-round student development at RDR School."},
  '/about/principal':{title:"Principal's Message | RDR School",description:"Read the historical RDR School principal's message about academic learning, character and school values."},
  '/about/philosophy':{title:'School Vision & Philosophy | RDR School',description:"Learn about RDR School's vision of broadening minds, independent thinking, compassion and growth."},
  '/academics':{title:'Academics | RDR School Haluhera',description:'Explore learning stages, academics and school activities at Rao Deena Ram Vidya Vihar in Rewari.'},
  '/campus':{title:'Campus & Facilities | RDR School',description:'Explore the school campus, science labs, computer lab, smart classrooms, sports and student facilities.'},
  '/admissions':{title:'Admissions & Enquiry | RDR School',description:'Contact RDR School about admissions, classes, fees, transport and visits. Submit a quick enquiry to the school office.'},
  '/gallery':{title:'School Life & Gallery | RDR School',description:'See moments of student life, culture, sports and school achievements from the RDR School gallery.'},
  '/mandatory-public-disclosure':{title:'CBSE Mandatory Public Disclosure | RDR School',description:'Find RDR School CBSE affiliation details, historic Appendix IX disclosures and school documents.'},
  '/resources':{title:'School Resources | RDR School',description:'Find school academic information, the calendar and available documents from RDR School.'},
  '/contact':{title:'Contact RDR School | Haluhera, Rewari',description:'Contact Rao Deena Ram Vidya Vihar School in Haluhera, Rewari, Haryana. Phone, email, directions and enquiries.'},
};

export { pageContent, pathPanel, pageMetadata };
