import { ArrowRight, ArrowUpRight, BookOpenText, CalendarDays, ExternalLink, FileText, GraduationCap, Mail, MapPin, Phone, ShieldCheck } from 'lucide-react';
import { disclosedDocuments, facilities, photographs, school } from '../../../../config/content';
import { PanelPhoto } from '../panelPhoto';

const AcademicPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <div className="rdr-sheet-overview">
    <div className="rdr-sheet-lead"><span className="rdr-sheet-kicker">LEARNING WITH PURPOSE</span><h3>Every stage. Every possibility.</h3>
      <p>Explore a learning environment that values curiosity, character and progress inside and outside the classroom.</p></div>
    <div className="rdr-sheet-grid">
      {[['Foundational years','A strong start through exploration and playful learning.'],['Primary','Literacy, numeracy, creativity and confidence.'],['Middle school','Independent thinking and deeper understanding.'],['Secondary & senior secondary','Preparation for future studies and opportunities.']].map(([title,description],i)=><article key={title}><span>0{i+1}</span><BookOpenText size={21}/><h4>{title}</h4><p>{description}</p></article>)}
    </div>
    <button type="button" className="rdr-sheet-action" onClick={onEnquire}>Ask about a class or curriculum <ArrowRight size={17}/></button>
  </div>
);

const CampusPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <div className="rdr-sheet-overview">
    <div className="rdr-sheet-lead"><span className="rdr-sheet-kicker">LEARNING SPACES</span><h3>Explore. Discover. Belong.</h3><p>A snapshot of the facilities described on RDR's original school website.</p></div>
    <div className="rdr-sheet-campus">
      <PanelPhoto index={2} alt="Children enjoying outdoor activities at RDR School"/>
      <div className="rdr-sheet-campus__list">
        {facilities.map((facility,i)=><div key={facility.title}><span>{String(i+1).padStart(2,'0')}</span><div><strong>{facility.title}</strong><small>{facility.description}</small></div></div>)}
      </div>
    </div>
    <button type="button" className="rdr-sheet-action" onClick={onEnquire}>Ask about a campus visit <ArrowRight size={17}/></button>
  </div>
);

const GalleryPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <div className="rdr-sheet-overview">
    <div className="rdr-sheet-lead"><span className="rdr-sheet-kicker">REAL MOMENTS · SCHOOL ARCHIVE</span><h3>Life at RDR, in pictures.</h3><p>Celebrations, achievements and experiences from the original school photo collection.</p></div>
    <div className="rdr-sheet-gallery">{photographs.slice(0,6).map((photo,i)=><figure key={photo.src} className={i===0?'featured':''}><PanelPhoto index={i} alt={photo.alt}/><figcaption>{photo.category}</figcaption></figure>)}</div>
    <button type="button" className="rdr-sheet-action" onClick={onEnquire}>Arrange a school visit <ArrowRight size={17}/></button>
  </div>
);

const AdmissionsPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <div className="rdr-sheet-overview">
    <div className="rdr-sheet-lead"><span className="rdr-sheet-kicker">WELCOME TO OUR SCHOOL</span><h3>Your next chapter begins here.</h3><p>Whether you are considering admission or simply exploring options, the school can help with classes, fees, transport and campus visits.</p></div>
    <div className="rdr-sheet-grid rdr-sheet-grid--three">
      {[[GraduationCap,'Explore classes','Ask about admission opportunities.'],[CalendarDays,'Visit the campus','Arrange a time to see the school.'],[Phone,'Speak to our team','Get answers to your questions.']].map(([Icon,title,description])=>{
        const Component=Icon as typeof GraduationCap;
        return <article key={title as string}><Component size={22}/><h4>{title as string}</h4><p>{description as string}</p></article>;
      })}
    </div>
    <button type="button" className="rdr-sheet-action" onClick={onEnquire}>Start a quick enquiry <ArrowRight size={17}/></button>
  </div>
);

const ResourcesPanel = () => (
  <div className="rdr-sheet-overview">
    <div className="rdr-sheet-lead"><span className="rdr-sheet-kicker">SCHOOL INFORMATION</span><h3>Useful links, less searching.</h3><p>Historical documents and useful school information from the original website.</p></div>
    <div className="rdr-sheet-resources">
      {disclosedDocuments.filter(item=>item.section==='C').map(item=><a href={item.url} target="_blank" rel="noreferrer" key={item.label}><FileText size={20}/><span>{item.label}</span><ExternalLink size={16}/></a>)}
      <a href={school.disclosurePdf} target="_blank" rel="noreferrer"><ShieldCheck size={20}/><span>2023 Mandatory Public Disclosure archive</span><ExternalLink size={16}/></a>
    </div>
    <p className="rdr-sheet-footnote">Document availability and validity are not confirmed for the current academic session.</p>
  </div>
);

const ContactPanel = ({onEnquire}:{onEnquire:()=>void}) => (
  <div className="rdr-sheet-overview">
    <div className="rdr-sheet-lead"><span className="rdr-sheet-kicker">WE'RE HERE TO HELP</span><h3>Reach RDR School.</h3><p>Contact the office or send your enquiry without navigating through multiple pages.</p></div>
    <div className="rdr-sheet-contact">
      <a href="tel:+919050077550"><Phone size={24}/><span><small>CALL OUR SCHOOL</small><strong>{school.phone}</strong></span><ArrowUpRight size={17}/></a>
      <a href={`mailto:${school.email}`}><Mail size={24}/><span><small>WRITE TO US</small><strong>{school.email}</strong></span><ArrowUpRight size={17}/></a>
      <a target="_blank" rel="noreferrer" href="https://www.google.com/maps/search/?api=1&query=Rao+Deena+Ram+Vidya+Vihar+Haluhera">
        <MapPin size={24}/><span><small>VISIT THE SCHOOL</small><strong>Haluhera, Rewari, Haryana</strong></span><ArrowUpRight size={17}/>
      </a>
    </div>
    <button type="button" className="rdr-sheet-action" onClick={onEnquire}>Enquire online <ArrowRight size={17}/></button>
  </div>
);

export { AcademicPanel, CampusPanel, GalleryPanel, AdmissionsPanel, ResourcesPanel, ContactPanel };
