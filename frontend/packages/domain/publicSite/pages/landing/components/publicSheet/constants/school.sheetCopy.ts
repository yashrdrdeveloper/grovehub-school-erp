import { GraduationCap, HeartHandshake, Landmark, Sparkles } from '@Components/icons';
import type { SchoolPanel } from '../types/schoolPanel';

const panelTitles: Record<SchoolPanel, string> = {
  enquiry:'Enquire at RDR School',
  about:'Discover RDR School',
  history:'Our school story',
  director:"Director's message",
  principal:"Principal's message",
  philosophy:'Our philosophy',
  academics:'Learning at RDR',
  campus:'Campus & facilities',
  gallery:'Life at RDR',
  admissions:'Admissions',
  resources:'School resources',
  contact:'Visit & contact',
};

const legacyStories = {
  history:{
    eyebrow:'ESTABLISHED IN APRIL 2006',
    title:'Rooted in Haluhera. Built for generations.',
    description:'The school was established in April 2006 by Shri Lal Singh Yadav (Retired HES-II), in memory of Late Shri Rao Deena Ram Ji. Its original mission was to bring quality educational opportunities closer to families around Rewari.',
    highlights:['Serving the local community since 2006','Growing with students and their families','Learning, character and responsibility'],
    link:'https://rdrschool.com/our-school.html',
    linkLabel:'Read the original school history',
  },
  director:{
    eyebrow:'SHRI LAL SINGH YADAV · LEGACY SCHOOL WEBSITE',
    title:'An education for the whole child.',
    description:"The director's published message highlights all-round development, school values and opportunities spanning academics, sports, arts, drama and co-curricular activities. It emphasizes preparing children for life, not simply examinations.",
    highlights:['Holistic development','Sports, arts and academics','A foundation for tomorrow'],
    link:'https://rdrschool.com/director-message.html',
    linkLabel:"Read the original director's message",
  },
  principal:{
    eyebrow:'SHRI ARVIND YADAV · LEGACY SCHOOL WEBSITE',
    title:'Learning with discipline and confidence.',
    description:"The principal's published message describes the school's development from modest beginnings and its commitment to quality education, experienced teachers, laboratories, co-curricular learning and character building.",
    highlights:['Quality education and guidance','Discipline and gratitude','Confidence beyond the classroom'],
    link:'https://rdrschool.com/Principal-message.html',
    linkLabel:"Read the original principal's message",
  },
  philosophy:{
    eyebrow:'OUR VISION · BROADENING MINDS',
    title:'Be a lamp unto yourself.',
    description:"RDR's original philosophy calls for a caring and creative environment that nurtures students' social, physical and intellectual growth. It places compassion, critical thinking and independence alongside academic achievement.",
    highlights:['Curiosity and independent thinking','Care, fairness and compassion','Prepared for a changing world'],
    link:'https://rdrschool.com/our-philosophy.html',
    linkLabel:'Read the school philosophy',
  },
} as const;

type StoryKey = keyof typeof legacyStories;
const storyNavigation: {id:StoryKey;label:string;icon:typeof Landmark}[] = [
  {id:'history',label:'School story',icon:Landmark},
  {id:'director',label:'Director',icon:HeartHandshake},
  {id:'principal',label:'Principal',icon:GraduationCap},
  {id:'philosophy',label:'Our vision',icon:Sparkles},
];
const overviewPanels: Record<SchoolPanel,string> = {
  enquiry:'A quick request, without a lengthy application.',
  about:'The people, values and journey behind RDR School.',
  history:'How our school began and grew.',
  director:'A message about every child’s growth.',
  principal:'A message about learning and character.',
  philosophy:'The values guiding our community.',
  academics:'Discover the learning journey.',
  campus:'See our learning spaces and facilities.',
  gallery:'Moments from our school life.',
  admissions:'Take the first step towards RDR.',
  resources:'Access useful school information.',
  contact:'We are happy to hear from you.',
};

export { panelTitles, legacyStories, storyNavigation, overviewPanels };
export type { StoryKey };
