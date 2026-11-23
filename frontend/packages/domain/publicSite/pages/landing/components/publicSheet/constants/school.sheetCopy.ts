import { GraduationCap, HeartHandshake, Landmark, Sparkles } from '@Components/icons';
import type { SchoolPanel } from '../types/schoolPanel';

const panelTitles: Record<SchoolPanel, string> = {
  enquiry: __('Enquire at RDR School'),
  about: __('Discover RDR School'),
  history: __('Our school story'),
  director: __("Director's message"),
  principal: __("Principal's message"),
  philosophy: __('Our philosophy'),
  academics: __('Learning at RDR'),
  campus: __('Campus & facilities'),
  gallery: __('Life at RDR'),
  admissions: __('Admissions'),
  resources: __('School resources'),
  contact: __('Visit & contact'),
};

const legacyStories = {
  history: {
    eyebrow: __('ESTABLISHED IN APRIL 2006'),
    title: __('Rooted in Haluhera. Built for generations.'),
    description:
      __('The school was established in April 2006 by Shri Lal Singh Yadav (Retired HES-II), in memory of Late Shri Rao Deena Ram Ji. Its original mission was to bring quality educational opportunities closer to families around Rewari.'),
    highlights: [
      __('Serving the local community since 2006'),
      __('Growing with students and their families'),
      __('Learning, character and responsibility'),
    ],
    link: 'https://rdrschool.com/our-school.html',
    linkLabel: __('Read the original school history'),
  },
  director: {
    eyebrow: __('SHRI LAL SINGH YADAV · LEGACY SCHOOL WEBSITE'),
    title: __('An education for the whole child.'),
    description:
      __("The director's published message highlights all-round development, school values and opportunities spanning academics, sports, arts, drama and co-curricular activities. It emphasizes preparing children for life, not simply examinations."),
    highlights: [__('Holistic development'), __('Sports, arts and academics'), __('A foundation for tomorrow')],
    link: 'https://rdrschool.com/director-message.html',
    linkLabel: __("Read the original director's message"),
  },
  principal: {
    eyebrow: __('SHRI ARVIND YADAV · LEGACY SCHOOL WEBSITE'),
    title: __('Learning with discipline and confidence.'),
    description:
      __("The principal's published message describes the school's development from modest beginnings and its commitment to quality education, experienced teachers, laboratories, co-curricular learning and character building."),
    highlights: [
      __('Quality education and guidance'),
      __('Discipline and gratitude'),
      __('Confidence beyond the classroom'),
    ],
    link: 'https://rdrschool.com/Principal-message.html',
    linkLabel: __("Read the original principal's message"),
  },
  philosophy: {
    eyebrow: __('OUR VISION · BROADENING MINDS'),
    title: __('Be a lamp unto yourself.'),
    description:
      __("RDR's original philosophy calls for a caring and creative environment that nurtures students' social, physical and intellectual growth. It places compassion, critical thinking and independence alongside academic achievement."),
    highlights: [
      __('Curiosity and independent thinking'),
      __('Care, fairness and compassion'),
      __('Prepared for a changing world'),
    ],
    link: 'https://rdrschool.com/our-philosophy.html',
    linkLabel: __('Read the school philosophy'),
  },
} as const;

type StoryKey = keyof typeof legacyStories;
const storyNavigation: { id: StoryKey; label: string; icon: typeof Landmark }[] = [
  { id: 'history', label: __('School story'), icon: Landmark },
  { id: 'director', label: __('Director'), icon: HeartHandshake },
  { id: 'principal', label: __('Principal'), icon: GraduationCap },
  { id: 'philosophy', label: __('Our vision'), icon: Sparkles },
];
const overviewPanels: Record<SchoolPanel, string> = {
  enquiry: __('A quick request, without a lengthy application.'),
  about: __('The people, values and journey behind RDR School.'),
  history: __('How our school began and grew.'),
  director: __('A message about every child’s growth.'),
  principal: __('A message about learning and character.'),
  philosophy: __('The values guiding our community.'),
  academics: __('Discover the learning journey.'),
  campus: __('See our learning spaces and facilities.'),
  gallery: __('Moments from our school life.'),
  admissions: __('Take the first step towards RDR.'),
  resources: __('Access useful school information.'),
  contact: __('We are happy to hear from you.'),
};

export { panelTitles, legacyStories, storyNavigation, overviewPanels };
export type { StoryKey };
