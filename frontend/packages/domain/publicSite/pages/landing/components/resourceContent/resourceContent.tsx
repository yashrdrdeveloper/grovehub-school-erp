import { _map } from '@Base/utilities/map';
import { _filter } from '@Base/utilities/filter';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Box } from '@Components/components/Box';
import { ArrowUpRight, BookOpen, CalendarDays, Clock3, Download, ExternalLink, FileCheck2, GraduationCap, Mail, MapPin, Phone, ShieldCheck } from '@Components/icons';
import { disclosedDocuments, historicDisclosure, school } from '../../constants/school.content';
import { Heading } from '../siteMedia';
import { IntroPage } from '../aboutContent';
import { EnquiryForm } from '../enquiryForm';

const DisclosurePage = () => (
  <>
    <IntroPage eyebrow="CBSE · APPENDIX IX" title="Mandatory Public Disclosure." description="Transparent access to school information, the affiliation reference and documents identified in RDR School's historical CBSE disclosure."/>
    <Box as="section" className="rdr-container rdr-disclosure-note"><ShieldCheck size={24}/><Box as="div"><Box as="strong">Historical disclosure — verification required</Box><Box as="p">The document available on the legacy website is dated 18 June 2023 (SARAS 4.0). The figures and linked certificates below are historical and must not be represented as verified current compliance. The current Appendix IX disclosure and active certificates must be supplied by the school.</Box></Box><ExternalAnchor href={school.disclosurePdf} target="_blank" rel="noreferrer">View original PDF <ExternalLink size={15}/></ExternalAnchor></Box>
    <Box as="section" className="rdr-section rdr-container rdr-disclosure">
      <Box as="div" className="rdr-disclosure__group"><Heading eyebrow="A · GENERAL INFORMATION" title="School information"/><Box as="div" className="rdr-disclosure__table">{_map([
        ['School',school.fullName],['Affiliation number',school.affiliation],['School code',school.schoolCode],['Address',school.address],['Principal (2023 record)',historicDisclosure.principal],['Principal qualification (2023 record)',historicDisclosure.qualification],['Email (2023 record)',school.email],
      ], ([k,v])=><Box as="div" key={k}><Box as="strong">{k}</Box><Box as="span">{v}</Box></Box>)}</Box></Box>
      <Box as="div" className="rdr-disclosure__group"><Heading eyebrow="B · DOCUMENTS AND INFORMATION" title="Certificates & registrations"/><DocumentList section="B"/></Box>
      <Box as="div" className="rdr-disclosure__group"><Heading eyebrow="C · RESULTS AND ACADEMICS" title="Academic documents"/><DocumentList section="C"/></Box>
      <Box as="div" className="rdr-disclosure__group"><Heading eyebrow="D · TEACHING STAFF" title="Teaching details (2023 archive)"/><Box as="div" className="rdr-disclosure__table">{_map([
        ['Total teachers','44'],['PGT','13'],['TGT','8'],['PRT','14'],['Special educators','1'],['Counsellor / wellness teacher','1'],
      ], ([k,v])=><Box as="div" key={k}><Box as="strong">{k}</Box><Box as="span">{v}</Box></Box>)}</Box><Box as="p" className="rdr-source-note">CBSE's current requirements include teachers' names and qualifications; updated, authorized staff information is still required.</Box></Box>
      <Box as="div" className="rdr-disclosure__group"><Heading eyebrow="E · INFRASTRUCTURE" title="Infrastructure (2023 archive)"/><Box as="div" className="rdr-disclosure__table">{_map([
        ['Campus area','8,000 sq. m'],['Classrooms','38'],['Laboratories (including computer labs)','6'],['Internet facility','Yes'],['Girls toilets','20'],['Boys toilets','20'],
      ], ([k,v])=><Box as="div" key={k}><Box as="strong">{k}</Box><Box as="span">{v}</Box></Box>)}</Box><ExternalAnchor className="rdr-inline-link" href={historicDisclosure.inspectionVideo} rel="noreferrer" target="_blank">Historical school inspection video <ExternalLink size={16}/></ExternalAnchor></Box>
      <Box as="div" className="rdr-disclosure__group"><Heading eyebrow="OFFICIAL VERIFICATION" title="Check the affiliation reference"/><ExternalAnchor href={school.cbseSaras} target="_blank" rel="noreferrer" className="rdr-btn rdr-btn--outline">CBSE SARAS affiliation listing <ExternalLink size={16}/></ExternalAnchor></Box>
    </Box>
  </>
);

const DocumentList = ({section}:{section:'B'|'C'}) => (
  <Box as="div" className="rdr-documents">{_map(_filter(disclosedDocuments, (doc) => doc.section === section), doc=><ExternalAnchor href={doc.url} key={doc.label} target="_blank" rel="noreferrer"><FileCheck2 size={19}/><Box as="span">{doc.label}</Box><Download size={17}/></ExternalAnchor>)}</Box>
);

const ResourcesPage = () => (
  <>
    <IntroPage eyebrow="PARENT RESOURCES" title="Useful information, together." description="Find school resources and archived documents in one place. Current academic-session resources will be published when confirmed."/>
    <Box as="section" className="rdr-section rdr-container"><Box as="div" className="rdr-feature-grid">{_map([
      {title:'Academic calendar',icon:CalendarDays,detail:'2023 archived academic calendar',href:'https://rdrschool.com/pdf/AcadmicCalender.pdf'},
      {title:'Fee information',icon:FileCheck2,detail:'Archived structure; confirm current fees',href:'https://rdrschool.com/pdf/Fees.pdf'},
      {title:'Board examination results',icon:GraduationCap,detail:'Historical result document',href:'https://rdrschool.com/pdf/Results.pdf'},
      {title:'School brochure',icon:BookOpen,detail:'Legacy brochure — historical material',href:'https://rdrschool.com/pdf/Brochure.pdf'},
    ], item=><ExternalAnchor href={item.href} target="_blank" rel="noreferrer" className="rdr-feature rdr-feature--link" key={item.title}><item.icon size={26}/><Box as="h3">{item.title}</Box><Box as="p">{item.detail}</Box><ExternalLink size={16}/></ExternalAnchor>)}</Box>
      <Box as="div" className="rdr-historic-hours"><Box as="div"><Clock3 size={24}/><Box as="h2">School timings — legacy information</Box><Box as="p">Published on the previous school website. Please confirm current term hours and holidays with the school office.</Box></Box><Box as="div"><Box as="strong">Summer</Box><Box as="span">Nursery · 8:00 AM–12:00 noon</Box><Box as="span">Classes I–XII · 8:00 AM–3:00 PM</Box></Box><Box as="div"><Box as="strong">Winter</Box><Box as="span">Nursery · 8:30 AM–12:30 PM</Box><Box as="span">Classes I–XII · 8:30 AM–3:30 PM</Box></Box></Box>
    </Box>
  </>
);

const ContactPage = () => (
  <>
    <IntroPage eyebrow="GET IN TOUCH" title="We'd love to hear from you." description="Have a question about admissions, school life, transport or academics? Contact the school office or send a quick enquiry."/>
    <Box as="section" className="rdr-section rdr-container rdr-admissions-layout">
      <Box as="div" className="rdr-contact-details"><Heading eyebrow="VISIT US" title="Find RDR School."/>
        <Box as="p"><MapPin size={20}/>{school.address}</Box><Box as="p"><Phone size={20}/><ExternalAnchor href="tel:+919050077550">{school.phone}</ExternalAnchor></Box><Box as="p"><Mail size={20}/><ExternalAnchor href={`mailto:${school.email}`}>{school.email}</ExternalAnchor></Box>
        <Box as="div" className="rdr-contact-map"><MapPin size={32}/><Box as="strong">Haluhera · Rewari</Box><ExternalAnchor href="https://www.google.com/maps/search/?api=1&query=Rao+Deena+Ram+Vidya+Vihar+Haluhera" target="_blank" rel="noreferrer">Open directions <ArrowUpRight size={16}/></ExternalAnchor></Box>
        <Box as="p" className="rdr-source-note">Contact details were taken from the legacy website and 2023 disclosure; confirm before public launch.</Box>
      </Box><EnquiryForm/>
    </Box>
  </>
);

export { DisclosurePage };
export { ResourcesPage };
export { ContactPage };
