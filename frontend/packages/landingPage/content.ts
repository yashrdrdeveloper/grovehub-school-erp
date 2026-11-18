/**
 * RDR School public-site content inventory.
 *
 * Historical content and images: https://rdrschool.com/
 * The 2023 SARAS disclosure is an ARCHIVE, not a current CBSE certification.
 * Verify staff names, document validity, address, contacts and session-specific
 * admissions information with the school before publishing a production site.
 */
export const school = {
  name: 'Rao Deena Ram Vidya Vihar',
  fullName: 'Rao Deena Ram Vidya Vihar Senior Secondary School',
  shortName: 'RDR School',
  address: 'Village Haluhera, P.O. Musepur, District Rewari, Haryana 123401',
  phone: '+91 90500 77550', // Existing website; verify before launch.
  email: 'rdrvvsshaluhera@gmail.com', // 2023 SARAS document; website differs.
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
    src: '/school-images/slide3.jpg', originalSrc: 'https://www.rdrschool.com/images/slide3.jpg',
    title: 'Learning through achievement',
    alt: 'RDR School students with awards and teachers at a school event',
    category: 'Achievements',
    position: 'center 42%',
  },
  {
    src: '/school-images/gallery-13.jpg', originalSrc: 'https://www.rdrschool.com/images/gallery/13.jpg',
    title: 'The joy of performing',
    alt: 'RDR School students performing a group cultural dance at a school celebration',
    category: 'Culture',
    position: 'center 38%',
  },
  {
    src: '/school-images/gallery-6.jpg', originalSrc: 'https://www.rdrschool.com/images/gallery/6.jpg',
    title: 'Everyday moments of discovery',
    alt: 'Young RDR School children playing outside at school',
    category: 'Campus Life',
    position: 'center center',
  },
  {
    src: '/school-images/gallery-10.jpg', originalSrc: 'https://www.rdrschool.com/images/gallery/10.jpg',
    title: 'Confidence in every move',
    alt: 'Students participating in a taekwondo demonstration at an RDR School event',
    category: 'Sports & Activities',
    position: 'center 35%',
  },
  {
    src: '/school-images/gallery-2.jpg', originalSrc: 'https://rdrschool.com/images/gallery/2.jpg',
    title: 'Growing together',
    alt: 'RDR School students smiling together in their school uniforms',
    category: 'Student Life',
    position: 'center center',
  },
  {
    src: '/school-images/slide2.jpg', originalSrc: 'https://www.rdrschool.com/images/slide2.jpg',
    title: 'Celebrating our traditions',
    alt: 'RDR School students celebrating cultural traditions in costume',
    category: 'Culture',
    position: 'center center',
  },
  {
    src: '/school-images/gallery-14.jpg', originalSrc: 'https://www.rdrschool.com/images/gallery/14.jpg',
    title: 'A stage for every talent',
    alt: 'Children dancing during an RDR School cultural event',
    category: 'Events',
    position: 'center center',
  },
  {
    src: '/school-images/gallery-7.jpg', originalSrc: 'https://www.rdrschool.com/images/gallery/7.jpg',
    title: 'Little moments, big milestones',
    alt: 'Students receiving awards on stage at an RDR School ceremony',
    category: 'Achievements',
    position: 'center center',
  },
] as const;

export const facilities = [
  { title:'Computer Lab', description:'Dedicated spaces for hands-on digital learning and computer education.', icon:'computer' },
  { title:'Science Laboratories', description:'Practical learning spaces for observation, experiments and discovery.', icon:'flask-conical' },
  { title:'Smart Classrooms', description:'Classroom learning supported by digital resources and teaching aids.', icon:'presentation' },
  { title:'Library', description:'Reading and reference materials that encourage independent learning.', icon:'book-open' },
  { title:'Sports & Activities', description:'Opportunities for physical education, teamwork and self-confidence.', icon:'trophy' },
  { title:'Transport', description:'School transportation for surrounding villages; contact the office for current routes.', icon:'bus' },
  { title:'Medical & Wellbeing', description:'School medical room and a focus on student care.', icon:'heart-pulse' },
  { title:'Campus Safety', description:'Facilities and safety measures described in historical school records; verify current certifications.', icon:'shield-check' },
] as const;

export const enquiryTypes = [
  { value:'ADMISSION', label:'Admission' },
  { value:'CAMPUS_VISIT', label:'Campus visit' },
  { value:'FEES', label:'Fees & scholarships' },
  { value:'TRANSPORT', label:'Transport' },
  { value:'ACADEMICS', label:'Academics' },
  { value:'GENERAL', label:'General enquiry' },
] as const;

export const disclosedDocuments = [
  { label:'CBSE affiliation / extension letter', url:'https://rdrschool.com/pdf/Affiliation.pdf', section:'B' },
  { label:'Society registration certificate', url:'https://rdrschool.com/pdf/Society.PDF', section:'B' },
  { label:'No Objection Certificate (NOC)', url:'https://rdrschool.com/pdf/NOC.pdf', section:'B' },
  { label:'Recognition certificate', url:'https://rdrschool.com/pdf/Recognition.pdf', section:'B' },
  { label:'Building safety certificate', url:'https://rdrschool.com/pdf/Building.PDF', section:'B' },
  { label:'Fire safety certificate', url:'https://rdrschool.com/pdf/FireSafety.pdf', section:'B' },
  { label:'Self-certification', url:'https://rdrschool.com/pdf/SelfDecleration.pdf', section:'B' },
  { label:'Water, health and sanitation', url:'https://rdrschool.com/pdf/Water.pdf', section:'B' },
  { label:'School fee structure', url:'https://rdrschool.com/pdf/Fees.pdf', section:'C' },
  { label:'Academic calendar', url:'https://rdrschool.com/pdf/AcadmicCalender.pdf', section:'C' },
  { label:'School Management Committee (SMC)', url:'https://rdrschool.com/pdf/SMC.pdf', section:'C' },
  { label:'Parent Teacher Association (PTA)', url:'https://rdrschool.com/pdf/PTA.pdf', section:'C' },
  { label:'Board examination results', url:'https://rdrschool.com/pdf/Results.pdf', section:'C' },
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
  internet: 'Yes',
  girlsToilets: 20,
  boysToilets: 20,
  inspectionVideo: 'https://youtu.be/pJrIfB_zZPo',
} as const;
