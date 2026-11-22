import { ArrowUpRight, Clock3, Download, ExternalLink, FileCheck2, Mail, MapPin, Phone, ShieldCheck } from 'lucide-react';
import { disclosedDocuments, historicDisclosure, school } from '../../config/content';
import { Heading } from '../siteMedia';
import { IntroPage } from '../aboutContent';
import { EnquiryForm } from '../enquiryForm';

const DisclosurePage = () => (
  <>
    <IntroPage eyebrow="CBSE · APPENDIX IX" title="Mandatory Public Disclosure." description="Transparent access to school information, the affiliation reference and documents identified in RDR School's historical CBSE disclosure."/>
    <section className="rdr-container rdr-disclosure-note"><ShieldCheck size={24}/><div><strong>Historical disclosure — verification required</strong><p>The document available on the legacy website is dated 18 June 2023 (SARAS 4.0). The figures and linked certificates below are historical and must not be represented as verified current compliance. The current Appendix IX disclosure and active certificates must be supplied by the school.</p></div><a href={school.disclosurePdf} target="_blank" rel="noreferrer">View original PDF <ExternalLink size={15}/></a></section>
    <section className="rdr-section rdr-container rdr-disclosure">
      <div className="rdr-disclosure__group"><Heading eyebrow="A · GENERAL INFORMATION" title="School information"/><div className="rdr-disclosure__table">{[
        ['School',school.fullName],['Affiliation number',school.affiliation],['School code',school.schoolCode],['Address',school.address],['Principal (2023 record)',historicDisclosure.principal],['Principal qualification (2023 record)',historicDisclosure.qualification],['Email (2023 record)',school.email],
      ].map(([k,v])=><div key={k}><strong>{k}</strong><span>{v}</span></div>)}</div></div>
      <div className="rdr-disclosure__group"><Heading eyebrow="B · DOCUMENTS AND INFORMATION" title="Certificates & registrations"/><DocumentList section="B"/></div>
      <div className="rdr-disclosure__group"><Heading eyebrow="C · RESULTS AND ACADEMICS" title="Academic documents"/><DocumentList section="C"/></div>
      <div className="rdr-disclosure__group"><Heading eyebrow="D · TEACHING STAFF" title="Teaching details (2023 archive)"/><div className="rdr-disclosure__table">{[
        ['Total teachers','44'],['PGT','13'],['TGT','8'],['PRT','14'],['Special educators','1'],['Counsellor / wellness teacher','1'],
      ].map(([k,v])=><div key={k}><strong>{k}</strong><span>{v}</span></div>)}</div><p className="rdr-source-note">CBSE's current requirements include teachers' names and qualifications; updated, authorized staff information is still required.</p></div>
      <div className="rdr-disclosure__group"><Heading eyebrow="E · INFRASTRUCTURE" title="Infrastructure (2023 archive)"/><div className="rdr-disclosure__table">{[
        ['Campus area','8,000 sq. m'],['Classrooms','38'],['Laboratories (including computer labs)','6'],['Internet facility','Yes'],['Girls toilets','20'],['Boys toilets','20'],
      ].map(([k,v])=><div key={k}><strong>{k}</strong><span>{v}</span></div>)}</div><a className="rdr-inline-link" href={historicDisclosure.inspectionVideo} rel="noreferrer" target="_blank">Historical school inspection video <ExternalLink size={16}/></a></div>
      <div className="rdr-disclosure__group"><Heading eyebrow="OFFICIAL VERIFICATION" title="Check the affiliation reference"/><a href={school.cbseSaras} target="_blank" rel="noreferrer" className="rdr-btn rdr-btn--outline">CBSE SARAS affiliation listing <ExternalLink size={16}/></a></div>
    </section>
  </>
);

const DocumentList = ({section}:{section:'B'|'C'}) => (
  <div className="rdr-documents">{disclosedDocuments.filter(doc=>doc.section===section).map(doc=><a href={doc.url} key={doc.label} target="_blank" rel="noreferrer"><FileCheck2 size={19}/><span>{doc.label}</span><Download size={17}/></a>)}</div>
);

const ResourcesPage = () => (
  <>
    <IntroPage eyebrow="PARENT RESOURCES" title="Useful information, together." description="Find school resources and archived documents in one place. Current academic-session resources will be published when confirmed."/>
    <section className="rdr-section rdr-container"><div className="rdr-feature-grid">{[
      {title:'Academic calendar',icon:CalendarDays,detail:'2023 archived academic calendar',href:'https://rdrschool.com/pdf/AcadmicCalender.pdf'},
      {title:'Fee information',icon:FileCheck2,detail:'Archived structure; confirm current fees',href:'https://rdrschool.com/pdf/Fees.pdf'},
      {title:'Board examination results',icon:GraduationCap,detail:'Historical result document',href:'https://rdrschool.com/pdf/Results.pdf'},
      {title:'School brochure',icon:BookOpen,detail:'Legacy brochure — historical material',href:'https://rdrschool.com/pdf/Brochure.pdf'},
    ].map(item=><a href={item.href} target="_blank" rel="noreferrer" className="rdr-feature rdr-feature--link" key={item.title}><item.icon size={26}/><h3>{item.title}</h3><p>{item.detail}</p><ExternalLink size={16}/></a>)}</div>
      <div className="rdr-historic-hours"><div><Clock3 size={24}/><h2>School timings — legacy information</h2><p>Published on the previous school website. Please confirm current term hours and holidays with the school office.</p></div><div><strong>Summer</strong><span>Nursery · 8:00 AM–12:00 noon</span><span>Classes I–XII · 8:00 AM–3:00 PM</span></div><div><strong>Winter</strong><span>Nursery · 8:30 AM–12:30 PM</span><span>Classes I–XII · 8:30 AM–3:30 PM</span></div></div>
    </section>
  </>
);

const ContactPage = () => (
  <>
    <IntroPage eyebrow="GET IN TOUCH" title="We'd love to hear from you." description="Have a question about admissions, school life, transport or academics? Contact the school office or send a quick enquiry."/>
    <section className="rdr-section rdr-container rdr-admissions-layout">
      <div className="rdr-contact-details"><Heading eyebrow="VISIT US" title="Find RDR School."/>
        <p><MapPin size={20}/>{school.address}</p><p><Phone size={20}/><a href="tel:+919050077550">{school.phone}</a></p><p><Mail size={20}/><a href={`mailto:${school.email}`}>{school.email}</a></p>
        <div className="rdr-contact-map"><MapPin size={32}/><strong>Haluhera · Rewari</strong><a href="https://www.google.com/maps/search/?api=1&query=Rao+Deena+Ram+Vidya+Vihar+Haluhera" target="_blank" rel="noreferrer">Open directions <ArrowUpRight size={16}/></a></div>
        <p className="rdr-source-note">Contact details were taken from the legacy website and 2023 disclosure; confirm before public launch.</p>
      </div><EnquiryForm/>
    </section>
  </>
);

export { DisclosurePage };
export { ResourcesPage };
export { ContactPage };
