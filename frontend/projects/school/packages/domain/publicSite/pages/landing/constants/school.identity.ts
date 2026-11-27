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
