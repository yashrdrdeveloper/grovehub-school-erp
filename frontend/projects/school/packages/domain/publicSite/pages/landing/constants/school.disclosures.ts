export const disclosedDocuments = [
  {
    label: __('CBSE affiliation / extension letter'),
    url: 'https://rdrschool.com/pdf/Affiliation.pdf',
    section: 'B',
  },
  {
    label: __('Society registration certificate'),
    url: 'https://rdrschool.com/pdf/Society.PDF',
    section: 'B',
  },
  {
    label: __('No Objection Certificate (NOC)'),
    url: 'https://rdrschool.com/pdf/NOC.pdf',
    section: 'B',
  },
  {
    label: __('Recognition certificate'),
    url: 'https://rdrschool.com/pdf/Recognition.pdf',
    section: 'B',
  },
  {
    label: __('Building safety certificate'),
    url: 'https://rdrschool.com/pdf/Building.PDF',
    section: 'B',
  },
  {
    label: __('Fire safety certificate'),
    url: 'https://rdrschool.com/pdf/FireSafety.pdf',
    section: 'B',
  },
  {
    label: __('Self-certification'),
    url: 'https://rdrschool.com/pdf/SelfDecleration.pdf',
    section: 'B',
  },
  {
    label: __('Water, health and sanitation'),
    url: 'https://rdrschool.com/pdf/Water.pdf',
    section: 'B',
  },
  { label: __('School fee structure'), url: 'https://rdrschool.com/pdf/Fees.pdf', section: 'C' },
  {
    label: __('Academic calendar'),
    url: 'https://rdrschool.com/pdf/AcadmicCalender.pdf',
    section: 'C',
  },
  {
    label: __('School Management Committee (SMC)'),
    url: 'https://rdrschool.com/pdf/SMC.pdf',
    section: 'C',
  },
  {
    label: __('Parent Teacher Association (PTA)'),
    url: 'https://rdrschool.com/pdf/PTA.pdf',
    section: 'C',
  },
  {
    label: __('Board examination results'),
    url: 'https://rdrschool.com/pdf/Results.pdf',
    section: 'C',
  },
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
