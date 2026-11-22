import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { ArrowRight, ArrowUpRight, BookOpenText, CalendarDays, ExternalLink, FileText, GraduationCap, Mail, MapPin, Phone, ShieldCheck } from '@Components/icons';
import { disclosedDocuments, facilities, photographs, school } from '../../../../constants/school.content';
import { PanelPhoto } from '../panelPhoto';
import { getInfoGalleryClassName } from './helpers/infoPanels.getGalleryClassName';

const AcademicPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as="div" className="rdr-sheet-overview">
    <Box as="div" className="rdr-sheet-lead"><Box as="span" className="rdr-sheet-kicker">LEARNING WITH PURPOSE</Box><Box as="h3">Every stage. Every possibility.</Box>
      <Box as="p">Explore a learning environment that values curiosity, character and progress inside and outside the classroom.</Box></Box>
    <Box as="div" className="rdr-sheet-grid">
      {[['Foundational years','A strong start through exploration and playful learning.'],['Primary','Literacy, numeracy, creativity and confidence.'],['Middle school','Independent thinking and deeper understanding.'],['Secondary & senior secondary','Preparation for future studies and opportunities.']].map(([title,description],i)=><Box as="article" key={title}><Box as="span">0{i+1}</Box><BookOpenText size={21}/><Box as="h4">{title}</Box><Box as="p">{description}</Box></Box>)}
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>Ask about a class or curriculum <ArrowRight size={17}/></Button>
  </Box>
);

const CampusPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as="div" className="rdr-sheet-overview">
    <Box as="div" className="rdr-sheet-lead"><Box as="span" className="rdr-sheet-kicker">LEARNING SPACES</Box><Box as="h3">Explore. Discover. Belong.</Box><Box as="p">A snapshot of the facilities described on RDR's original school website.</Box></Box>
    <Box as="div" className="rdr-sheet-campus">
      <PanelPhoto index={2} alt="Children enjoying outdoor activities at RDR School"/>
      <Box as="div" className="rdr-sheet-campus__list">
        {facilities.map((facility,i)=><Box as="div" key={facility.title}><Box as="span">{String(i+1).padStart(2,'0')}</Box><Box as="div"><Box as="strong">{facility.title}</Box><Box as="small">{facility.description}</Box></Box></Box>)}
      </Box>
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>Ask about a campus visit <ArrowRight size={17}/></Button>
  </Box>
);

const GalleryPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as="div" className="rdr-sheet-overview">
    <Box as="div" className="rdr-sheet-lead"><Box as="span" className="rdr-sheet-kicker">REAL MOMENTS · SCHOOL ARCHIVE</Box><Box as="h3">Life at RDR, in pictures.</Box><Box as="p">Celebrations, achievements and experiences from the original school photo collection.</Box></Box>
    <Box as="div" className="rdr-sheet-gallery">{photographs.slice(0,6).map((photo,i)=><figure key={photo.src} className={getInfoGalleryClassName(i)}><PanelPhoto index={i} alt={photo.alt}/><figcaption>{photo.category}</figcaption></figure>)}</Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>Arrange a school visit <ArrowRight size={17}/></Button>
  </Box>
);

const AdmissionsPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as="div" className="rdr-sheet-overview">
    <Box as="div" className="rdr-sheet-lead"><Box as="span" className="rdr-sheet-kicker">WELCOME TO OUR SCHOOL</Box><Box as="h3">Your next chapter begins here.</Box><Box as="p">Whether you are considering admission or simply exploring options, the school can help with classes, fees, transport and campus visits.</Box></Box>
    <Box as="div" className="rdr-sheet-grid rdr-sheet-grid--three">
      {[[GraduationCap,'Explore classes','Ask about admission opportunities.'],[CalendarDays,'Visit the campus','Arrange a time to see the school.'],[Phone,'Speak to our team','Get answers to your questions.']].map(([Icon,title,description])=>{
        const Component=Icon as typeof GraduationCap;
        return <Box as="article" key={title as string}><Component size={22}/><Box as="h4">{title as string}</Box><Box as="p">{description as string}</Box></Box>;
      })}
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>Start a quick enquiry <ArrowRight size={17}/></Button>
  </Box>
);

const ResourcesPanel = () => (
  <Box as="div" className="rdr-sheet-overview">
    <Box as="div" className="rdr-sheet-lead"><Box as="span" className="rdr-sheet-kicker">SCHOOL INFORMATION</Box><Box as="h3">Useful links, less searching.</Box><Box as="p">Historical documents and useful school information from the original website.</Box></Box>
    <Box as="div" className="rdr-sheet-resources">
      {disclosedDocuments.filter(item=>item.section==='C').map(item=><ExternalAnchor href={item.url} target="_blank" rel="noreferrer" key={item.label}><FileText size={20}/><Box as="span">{item.label}</Box><ExternalLink size={16}/></ExternalAnchor>)}
      <ExternalAnchor href={school.disclosurePdf} target="_blank" rel="noreferrer"><ShieldCheck size={20}/><Box as="span">2023 Mandatory Public Disclosure archive</Box><ExternalLink size={16}/></ExternalAnchor>
    </Box>
    <Box as="p" className="rdr-sheet-footnote">Document availability and validity are not confirmed for the current academic session.</Box>
  </Box>
);

const ContactPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as="div" className="rdr-sheet-overview">
    <Box as="div" className="rdr-sheet-lead"><Box as="span" className="rdr-sheet-kicker">WE'RE HERE TO HELP</Box><Box as="h3">Reach RDR School.</Box><Box as="p">Contact the office or send your enquiry without navigating through multiple pages.</Box></Box>
    <Box as="div" className="rdr-sheet-contact">
      <ExternalAnchor href="tel:+919050077550"><Phone size={24}/><Box as="span"><Box as="small">CALL OUR SCHOOL</Box><Box as="strong">{school.phone}</Box></Box><ArrowUpRight size={17}/></ExternalAnchor>
      <ExternalAnchor href={`mailto:${school.email}`}><Mail size={24}/><Box as="span"><Box as="small">WRITE TO US</Box><Box as="strong">{school.email}</Box></Box><ArrowUpRight size={17}/></ExternalAnchor>
      <ExternalAnchor target="_blank" rel="noreferrer" href="https://www.google.com/maps/search/?api=1&query=Rao+Deena+Ram+Vidya+Vihar+Haluhera">
        <MapPin size={24}/><Box as="span"><Box as="small">VISIT THE SCHOOL</Box><Box as="strong">Haluhera, Rewari, Haryana</Box></Box><ArrowUpRight size={17}/>
      </ExternalAnchor>
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>Enquire online <ArrowRight size={17}/></Button>
  </Box>
);

export { AcademicPanel, CampusPanel, GalleryPanel, AdmissionsPanel, ResourcesPanel, ContactPanel };
