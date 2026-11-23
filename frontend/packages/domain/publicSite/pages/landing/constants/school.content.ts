/**
 * RDR School public-site content inventory.
 *
 * Historical content and images: https://rdrschool.com/
 * The 2023 SARAS disclosure is an ARCHIVE, not a current CBSE certification.
 * Verify staff names, document validity, address, contacts and session-specific
 * admissions information with the school before publishing a production site.
 */
export const school = {
  name: __('Rao Deena Ram Vidya Vihar'),
  fullName: __('Rao Deena Ram Vidya Vihar Senior Secondary School'),
  shortName: __('RDR School'),
  address: __('Village Haluhera, P.O. Musepur, District Rewari, Haryana 123401'),
  phone: '+91 90500 77550', // Existing website; verify before launch.
  email: 'rdrvvssshaluhera@gmail.com', // 2023 SARAS document; website differs.
  affiliation: '531923',
  schoolCode: '41925',
  founded: '2006',
  website: 'https://rdrschool.com',
  disclosurePdf: 'https://rdrschool.com/pdf/Mandatory.pdf',
  cbseSaras: 'https://saras.cbse.gov.in/SARAS/AffiliatedList/AfflicationDetails/531923',
} as const;

/** Original RDR photo references, not stock images. Run `pnpm assets:school`
 * to copy the school's originals into public/school-images. If an image
 * has not been copied yet, the UI falls back to the verified legacy URL.
 */
export const photographs = [
  {
    src: '/school-images/slide3.jpg', originalSrc: 'https://rdrschool.com/images/slide3.jpg',
    title: __('Learning through achievement'),
    alt: __('RDR School students with awards and teachers at a school event'),
    category: __('Achievements'),
    position: 'center 42%',
  },
  {
    src: '/school-images/gallery-13.jpg', originalSrc: 'https://rdrschool.com/images/gallery/13.jpg',
    title: __('The joy of performing'),
    alt: __('RDR School students performing a group cultural dance at a school celebration'),
    category: __('Culture'),
    position: 'center 38%',
  },
  {
    src: '/school-images/gallery-6.jpg', originalSrc: 'https://rdrschool.com/images/gallery/6.jpg',
    title: __('Everyday moments of discovery'),
    alt: __('Young RDR School children playing outside at school'),
    category: __('Campus Life'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-10.jpg', originalSrc: 'https://rdrschool.com/images/gallery/10.jpg',
    title: __('Confidence in every move'),
    alt: __('Students participating in a taekwondo demonstration at an RDR School event'),
    category: __('Sports & Activities'),
    position: 'center 35%',
  },
  {
    src: '/school-images/slide5.jpg', originalSrc: 'https://rdrschool.com/images/slide5.jpg',
    title: __('Growing together'),
    alt: __('RDR School students smiling together in their school uniforms'),
    category: __('Student Life'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-11.jpg', originalSrc: 'https://rdrschool.com/images/gallery/11.jpg',
    title: __('Celebrating our traditions'),
    alt: __('RDR School students celebrating cultural traditions in costume'),
    category: __('Culture'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-4.jpg', originalSrc: 'https://rdrschool.com/images/gallery/4.jpg',
    title: __('A stage for every talent'),
    alt: __('Children dancing during an RDR School cultural event'),
    category: __('Events'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-7.jpg', originalSrc: 'https://rdrschool.com/images/gallery/7.jpg',
    title: __('Little moments, big milestones'),
    alt: __('Students receiving awards on stage at an RDR School ceremony'),
    category: __('Achievements'),
    position: 'center center',
  },
] as const;

export const facilities = [
  { title:__('Computer Lab'), description:__('Dedicated spaces for hands-on digital learning and computer education.'), icon:'computer' },
  { title:__('Science Laboratories'), description:__('Practical learning spaces for observation, experiments and discovery.'), icon:'flask-conical' },
  { title:__('Smart Classrooms'), description:__('Classroom learning supported by digital resources and teaching aids.'), icon:'presentation' },
  { title:__('Library'), description:__('Reading and reference materials that encourage independent learning.'), icon:'book-open' },
  { title:__('Sports & Activities'), description:__('Opportunities for physical education, teamwork and self-confidence.'), icon:'trophy' },
  { title:__('Transport'), description:__('School transportation for surrounding villages; contact the office for current routes.'), icon:'bus' },
  { title:__('Medical & Wellbeing'), description:__('School medical room and a focus on student care.'), icon:'heart-pulse' },
  { title:__('Campus Safety'), description:__('Facilities and safety measures described in historical school records; verify current certifications.'), icon:'shield-check' },
] as const;

export const enquiryTypes = [
  { value:'ADMISSION', label:__('Admission') },
  { value:'CAMPUS_VISIT', label:__('Campus visit') },
  { value:'FEES', label:__('Fees & scholarships') },
  { value:'TRANSPORT', label:__('Transport') },
  { value:'ACADEMICS', label:__('Academics') },
  { value:'GENERAL', label:__('General enquiry') },
] as const;

export const disclosedDocuments = [
  { label:__('CBSE affiliation / extension letter'), url:'https://rdrschool.com/pdf/Affiliation.pdf', section:'B' },
  { label:__('Society registration certificate'), url:'https://rdrschool.com/pdf/Society.PDF', section:'B' },
  { label:__('No Objection Certificate (NOC)'), url:'https://rdrschool.com/pdf/NOC.pdf', section:'B' },
  { label:__('Recognition certificate'), url:'https://rdrschool.com/pdf/Recognition.pdf', section:'B' },
  { label:__('Building safety certificate'), url:'https://rdrschool.com/pdf/Building.PDF', section:'B' },
  { label:__('Fire safety certificate'), url:'https://rdrschool.com/pdf/FireSafety.pdf', section:'B' },
  { label:__('Self-certification'), url:'https://rdrschool.com/pdf/SelfDecleration.pdf', section:'B' },
  { label:__('Water, health and sanitation'), url:'https://rdrschool.com/pdf/Water.pdf', section:'B' },
  { label:__('School fee structure'), url:'https://rdrschool.com/pdf/Fees.pdf', section:'C' },
  { label:__('Academic calendar'), url:'https://rdrschool.com/pdf/AcadmicCalender.pdf', section:'C' },
  { label:__('School Management Committee (SMC)'), url:'https://rdrschool.com/pdf/SMC.pdf', section:'C' },
  { label:__('Parent Teacher Association (PTA)'), url:'https://rdrschool.com/pdf/PTA.pdf', section:'C' },
  { label:__('Board examination results'), url:'https://rdrschool.com/pdf/Results.pdf', section:'C' },
] as const;

export const historicDisclosure = {
  principal: 'Arvind Yadav', // Name from 2023 Appendix IX only, not confirmed current.
  qualification: 'M.Sc., B.Ed.',
  teachers: 44,
  pgt: 13,
  tgt: 8,
  prt: 14,
  specialEducators: 1,
  counsellors: 1,
  campusSquareMetres: 8000,
  classrooms: 38,
  laboratories: 6,
  internet: __('Yes'),
  girlsToilets: 20,
  boysToilets: 20,
  inspectionVideo: 'https://youtu.be/pJrIfB_zZPo',
} as const;
