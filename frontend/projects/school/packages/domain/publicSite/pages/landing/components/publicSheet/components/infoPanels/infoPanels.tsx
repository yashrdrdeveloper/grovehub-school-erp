import { _map } from '@Base/utilities/map';
import { _filter } from '@Base/utilities/filter';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Button } from '@Components/components/Button';
import { Box , BOX_TAG} from '@Components/components/Box';
import { ArrowRight, ArrowUpRight, BookOpenText, CalendarDays, ExternalLink, FileText, GraduationCap, Mail, MapPin, Phone, ShieldCheck } from '@Components/icons';
import { disclosedDocuments, facilities, photographs, school } from '../../../../constants/school.content';
import { PanelPhoto } from '../panelPhoto';
import { getInfoGalleryClassName } from './helpers/infoPanels.getGalleryClassName';

const AcademicPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as={BOX_TAG.DIV} className="rdr-sheet-overview">
    <Box as={BOX_TAG.DIV} className="rdr-sheet-lead"><Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">{__("LEARNING WITH PURPOSE")}</Box><Box as={BOX_TAG.H3}>{__("Every stage. Every possibility.")}</Box>
      <Box as={BOX_TAG.P}>{__("Explore a learning environment that values curiosity, character and progress inside and outside the classroom.")}</Box></Box>
    <Box as={BOX_TAG.DIV} className="rdr-sheet-grid">
      {_map([['Foundational years','A strong start through exploration and playful learning.'],['Primary','Literacy, numeracy, creativity and confidence.'],['Middle school','Independent thinking and deeper understanding.'],['Secondary & senior secondary','Preparation for future studies and opportunities.']], ([title,description],i)=><Box as={BOX_TAG.ARTICLE} key={title}><Box as={BOX_TAG.SPAN}>0{i+1}</Box><BookOpenText size={21}/><Box as={BOX_TAG.H4}>{title}</Box><Box as={BOX_TAG.P}>{description}</Box></Box>)}
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>{__("Ask about a class or curriculum ")}<ArrowRight size={17}/></Button>
  </Box>
);

const CampusPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as={BOX_TAG.DIV} className="rdr-sheet-overview">
    <Box as={BOX_TAG.DIV} className="rdr-sheet-lead"><Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">{__("LEARNING SPACES")}</Box><Box as={BOX_TAG.H3}>{__("Explore. Discover. Belong.")}</Box><Box as={BOX_TAG.P}>{__("A snapshot of the facilities described on RDR's original school website.")}</Box></Box>
    <Box as={BOX_TAG.DIV} className="rdr-sheet-campus">
      <PanelPhoto index={2} alt={__("Children enjoying outdoor activities at RDR School")}/>
      <Box as={BOX_TAG.DIV} className="rdr-sheet-campus__list">
        {_map(facilities, (facility,i)=><Box as={BOX_TAG.DIV} key={facility.title}><Box as={BOX_TAG.SPAN}>{String(i+1).padStart(2,'0')}</Box><Box as={BOX_TAG.DIV}><Box as={BOX_TAG.STRONG}>{facility.title}</Box><Box as={BOX_TAG.SMALL}>{facility.description}</Box></Box></Box>)}
      </Box>
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>{__("Ask about a campus visit ")}<ArrowRight size={17}/></Button>
  </Box>
);

const GalleryPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as={BOX_TAG.DIV} className="rdr-sheet-overview">
    <Box as={BOX_TAG.DIV} className="rdr-sheet-lead"><Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">{__("REAL MOMENTS · SCHOOL ARCHIVE")}</Box><Box as={BOX_TAG.H3}>{__("Life at RDR, in pictures.")}</Box><Box as={BOX_TAG.P}>{__("Celebrations, achievements and experiences from the original school photo collection.")}</Box></Box>
    <Box as={BOX_TAG.DIV} className="rdr-sheet-gallery">{_map(photographs.slice(0,6), (photo,i)=><figure key={photo.src} className={getInfoGalleryClassName(i)}><PanelPhoto index={i} alt={photo.alt}/><figcaption>{photo.category}</figcaption></figure>)}</Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>{__("Arrange a school visit ")}<ArrowRight size={17}/></Button>
  </Box>
);

const AdmissionsPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as={BOX_TAG.DIV} className="rdr-sheet-overview">
    <Box as={BOX_TAG.DIV} className="rdr-sheet-lead"><Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">{__("WELCOME TO OUR SCHOOL")}</Box><Box as={BOX_TAG.H3}>{__("Your next chapter begins here.")}</Box><Box as={BOX_TAG.P}>{__("Whether you are considering admission or simply exploring options, the school can help with classes, fees, transport and campus visits.")}</Box></Box>
    <Box as={BOX_TAG.DIV} className="rdr-sheet-grid rdr-sheet-grid--three">
      {_map([[GraduationCap,'Explore classes','Ask about admission opportunities.'],[CalendarDays,'Visit the campus','Arrange a time to see the school.'],[Phone,'Speak to our team','Get answers to your questions.']], ([Icon,title,description])=>{
        const Component=Icon as typeof GraduationCap;
        return <Box as={BOX_TAG.ARTICLE} key={title as string}><Component size={22}/><Box as={BOX_TAG.H4}>{title as string}</Box><Box as={BOX_TAG.P}>{description as string}</Box></Box>;
      })}
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>{__("Start a quick enquiry ")}<ArrowRight size={17}/></Button>
  </Box>
);

const ResourcesPanel = () => (
  <Box as={BOX_TAG.DIV} className="rdr-sheet-overview">
    <Box as={BOX_TAG.DIV} className="rdr-sheet-lead"><Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">{__("SCHOOL INFORMATION")}</Box><Box as={BOX_TAG.H3}>{__("Useful links, less searching.")}</Box><Box as={BOX_TAG.P}>{__("Historical documents and useful school information from the original website.")}</Box></Box>
    <Box as={BOX_TAG.DIV} className="rdr-sheet-resources">
      {_map(_filter(disclosedDocuments, (item) => item.section === 'C'), item=><ExternalAnchor href={item.url} target="_blank" rel="noreferrer" key={item.label}><FileText size={20}/><Box as={BOX_TAG.SPAN}>{item.label}</Box><ExternalLink size={16}/></ExternalAnchor>)}
      <ExternalAnchor href={school.disclosurePdf} target="_blank" rel="noreferrer"><ShieldCheck size={20}/><Box as={BOX_TAG.SPAN}>{__("2023 Mandatory Public Disclosure archive")}</Box><ExternalLink size={16}/></ExternalAnchor>
    </Box>
    <Box as={BOX_TAG.P} className="rdr-sheet-footnote">{__("Document availability and validity are not confirmed for the current academic session.")}</Box>
  </Box>
);

const ContactPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <Box as={BOX_TAG.DIV} className="rdr-sheet-overview">
    <Box as={BOX_TAG.DIV} className="rdr-sheet-lead"><Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">{__("WE'RE HERE TO HELP")}</Box><Box as={BOX_TAG.H3}>{__("Reach RDR School.")}</Box><Box as={BOX_TAG.P}>{__("Contact the office or send your enquiry without navigating through multiple pages.")}</Box></Box>
    <Box as={BOX_TAG.DIV} className="rdr-sheet-contact">
      <ExternalAnchor href="tel:+919050077550"><Phone size={24}/><Box as={BOX_TAG.SPAN}><Box as={BOX_TAG.SMALL}>{__("CALL OUR SCHOOL")}</Box><Box as={BOX_TAG.STRONG}>{school.phone}</Box></Box><ArrowUpRight size={17}/></ExternalAnchor>
      <ExternalAnchor href={`mailto:${school.email}`}><Mail size={24}/><Box as={BOX_TAG.SPAN}><Box as={BOX_TAG.SMALL}>{__("WRITE TO US")}</Box><Box as={BOX_TAG.STRONG}>{school.email}</Box></Box><ArrowUpRight size={17}/></ExternalAnchor>
      <ExternalAnchor target="_blank" rel="noreferrer" href="https://www.google.com/maps/search/?api=1&query=Rao+Deena+Ram+Vidya+Vihar+Haluhera">
        <MapPin size={24}/><Box as={BOX_TAG.SPAN}><Box as={BOX_TAG.SMALL}>{__("VISIT THE SCHOOL")}</Box><Box as={BOX_TAG.STRONG}>{__("Haluhera, Rewari, Haryana")}</Box></Box><ArrowUpRight size={17}/>
      </ExternalAnchor>
    </Box>
    <Button unstyled type="button" className="rdr-sheet-action" onClick={onEnquire}>{__("Enquire online ")}<ArrowRight size={17}/></Button>
  </Box>
);

export { AcademicPanel, CampusPanel, GalleryPanel, AdmissionsPanel, ResourcesPanel, ContactPanel };
