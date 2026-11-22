import { useState, type ReactNode } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import {
  ArrowRight, ArrowUpRight, BookOpenText, CalendarDays, Check,
  ExternalLink, FileText, GraduationCap, HeartHandshake,
  Landmark, Mail, MapPin, Phone, ShieldCheck, Sparkles, X,
} from 'lucide-react';

import { disclosedDocuments, facilities, photographs, school } from '../../config/content';

export type SchoolPanel =
  | 'enquiry' | 'about' | 'history' | 'director' | 'principal' | 'philosophy'
  | 'academics' | 'campus' | 'gallery' | 'admissions' | 'resources' | 'contact';

const panelTitles: Record<SchoolPanel, string> = {
  enquiry:'Enquire at RDR School',
  about:'Discover RDR School',
  history:'Our school story',
  director:"Director's message",
  principal:"Principal's message",
  philosophy:'Our philosophy',
  academics:'Learning at RDR',
  campus:'Campus & facilities',
  gallery:'Life at RDR',
  admissions:'Admissions',
  resources:'School resources',
  contact:'Visit & contact',
};

const legacyStories = {
  history:{
    eyebrow:'ESTABLISHED IN APRIL 2006',
    title:'Rooted in Haluhera. Built for generations.',
    description:'The school was established in April 2006 by Shri Lal Singh Yadav (Retired HES-II), in memory of Late Shri Rao Deena Ram Ji. Its original mission was to bring quality educational opportunities closer to families around Rewari.',
    highlights:['Serving the local community since 2006','Growing with students and their families','Learning, character and responsibility'],
    link:'https://rdrschool.com/our-school.html',
    linkLabel:'Read the original school history',
  },
  director:{
    eyebrow:'SHRI LAL SINGH YADAV · LEGACY SCHOOL WEBSITE',
    title:'An education for the whole child.',
    description:"The director's published message highlights all-round development, school values and opportunities spanning academics, sports, arts, drama and co-curricular activities. It emphasizes preparing children for life, not simply examinations.",
    highlights:['Holistic development','Sports, arts and academics','A foundation for tomorrow'],
    link:'https://rdrschool.com/director-message.html',
    linkLabel:"Read the original director's message",
  },
  principal:{
    eyebrow:'SHRI ARVIND YADAV · LEGACY SCHOOL WEBSITE',
    title:'Learning with discipline and confidence.',
    description:"The principal's published message describes the school's development from modest beginnings and its commitment to quality education, experienced teachers, laboratories, co-curricular learning and character building.",
    highlights:['Quality education and guidance','Discipline and gratitude','Confidence beyond the classroom'],
    link:'https://rdrschool.com/Principal-message.html',
    linkLabel:"Read the original principal's message",
  },
  philosophy:{
    eyebrow:'OUR VISION · BROADENING MINDS',
    title:'Be a lamp unto yourself.',
    description:"RDR's original philosophy calls for a caring and creative environment that nurtures students' social, physical and intellectual growth. It places compassion, critical thinking and independence alongside academic achievement.",
    highlights:['Curiosity and independent thinking','Care, fairness and compassion','Prepared for a changing world'],
    link:'https://rdrschool.com/our-philosophy.html',
    linkLabel:'Read the school philosophy',
  },
} as const;

type StoryKey = keyof typeof legacyStories;
const storyNavigation: {id:StoryKey;label:string;icon:typeof Landmark}[] = [
  {id:'history',label:'School story',icon:Landmark},
  {id:'director',label:'Director',icon:HeartHandshake},
  {id:'principal',label:'Principal',icon:GraduationCap},
  {id:'philosophy',label:'Our vision',icon:Sparkles},
];

const PanelPhoto = ({index,alt}:{index:number;alt:string}) => {
  const photo=photographs[index];
  if (!photo) return null;
  return <div className="rdr-sheet-photo">
    <img loading="lazy" src={photo.src} alt={alt} style={{objectPosition:photo.position}}
      onError={event=>{
        if(event.currentTarget.dataset.retry) {event.currentTarget.hidden=true;return;}
        event.currentTarget.dataset.retry='1';
        event.currentTarget.src=photo.originalSrc;
      }}/>
    <span aria-hidden="true"><Sparkles size={17}/> RDR SCHOOL</span>
  </div>;
};

const StoryPanel = ({initial='history',navigate}:{initial?:StoryKey;navigate:(path:string)=>void}) => {
  const [selected,setSelected]=useState<StoryKey>(initial);
  const story=legacyStories[selected];
  return <div className="rdr-sheet-story">
    <div className="rdr-sheet-story__tabs" role="tablist" aria-label="School story topics">
      {storyNavigation.map(item=>{
        const Icon=item.icon;
        return <button type="button" role="tab" aria-selected={selected===item.id}
          key={item.id} className={selected===item.id?'active':''} onClick={()=>setSelected(item.id)}>
          <Icon size={16}/>{item.label}
        </button>;
      })}
    </div>
    <div className="rdr-sheet-story__layout" role="tabpanel">
      <div className="rdr-sheet-story__copy">
        <span className="rdr-sheet-kicker">{story.eyebrow}</span>
        <h3>{story.title}</h3>
        <p>{story.description}</p>
        <div className="rdr-sheet-story__points">
          {story.highlights.map(label=><span key={label}><Check size={15}/>{label}</span>)}
        </div>
        <div className="rdr-sheet-story__footer">
          <a href={story.link} target="_blank" rel="noreferrer">{story.linkLabel}<ExternalLink size={15}/></a>
          <button type="button" onClick={()=>navigate('enquiry')}>Ask us a question <ArrowRight size={16}/></button>
        </div>
        {(selected==='principal'||selected==='director')&&<small>Role and message are from the school's historical website. Please confirm current leadership before official publication.</small>}
      </div>
      <PanelPhoto index={selected==='history'?4:selected==='philosophy'?2:0}
        alt="Students participating in activities at RDR School"/>
    </div>
  </div>;
};

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

const overviewPanels: Record<SchoolPanel,string> = {
  enquiry:'A quick request, without a lengthy application.',
  about:'The people, values and journey behind RDR School.',
  history:'How our school began and grew.',
  director:'A message about every child’s growth.',
  principal:'A message about learning and character.',
  philosophy:'The values guiding our community.',
  academics:'Discover the learning journey.',
  campus:'See our learning spaces and facilities.',
  gallery:'Moments from our school life.',
  admissions:'Take the first step towards RDR.',
  resources:'Access useful school information.',
  contact:'We are happy to hear from you.',
};

export const PublicSheet = ({
  panel, onClose, onNavigate, enquiry,
}: {
  panel:SchoolPanel|null;
  onClose:()=>void;
  onNavigate:(panel:SchoolPanel)=>void;
  enquiry:ReactNode;
}) => {
  const selectedStory:StoryKey = panel==='director'||panel==='principal'||panel==='philosophy'||panel==='history'
    ?panel:'history';
  const story=panel==='about'||panel==='history'||panel==='director'||panel==='principal'||panel==='philosophy';
  const onEnquire=()=>onNavigate('enquiry');
  const title=panel?panelTitles[panel]:'RDR School information';
  return <Dialog.Root open={Boolean(panel)} onOpenChange={open=>{if(!open)onClose();}}>
    <Dialog.Portal>
      <Dialog.Overlay className="rdr-sheet-backdrop"/>
      <Dialog.Content className="rdr-sheet" aria-label={title}>
        <div className="rdr-sheet__handle" aria-hidden="true"/>
        <div className="rdr-sheet__header">
          <div><span className="rdr-sheet-kicker">RAO DEENA RAM VIDYA VIHAR</span><Dialog.Title>{title}</Dialog.Title>
            <Dialog.Description>{panel?overviewPanels[panel]:''}</Dialog.Description>
          </div>
          <Dialog.Close className="rdr-sheet__close" aria-label="Close and return to the page"><X size={22}/></Dialog.Close>
        </div>
        <div className="rdr-sheet__body" key={panel??'empty'}>
          {panel==='enquiry'&&<div className="rdr-sheet-enquiry"><div className="rdr-sheet-enquiry__intro"><GraduationCap size={26}/><strong>Let's talk about your child's future.</strong><p>Share your details. Our school office is only a message away.</p>
            <a href="tel:+919050077550"><Phone size={15}/> Prefer to call? {school.phone}</a></div>{enquiry}</div>}
          {story&&<StoryPanel initial={selectedStory} navigate={p=>onNavigate(p as SchoolPanel)}/>}
          {panel==='academics'&&<AcademicPanel onEnquire={onEnquire}/>}
          {panel==='campus'&&<CampusPanel onEnquire={onEnquire}/>}
          {panel==='gallery'&&<GalleryPanel onEnquire={onEnquire}/>}
          {panel==='admissions'&&<AdmissionsPanel onEnquire={onEnquire}/>}
          {panel==='resources'&&<ResourcesPanel/>}
          {panel==='contact'&&<ContactPanel onEnquire={onEnquire}/>}
        </div>
        {panel!=='enquiry'&&<div className="rdr-sheet__bottom"><span><ShieldCheck size={15}/> RDR School · Haluhera, Rewari</span><button onClick={onEnquire} type="button">Quick enquiry <ArrowRight size={16}/></button></div>}
      </Dialog.Content>
    </Dialog.Portal>
  </Dialog.Root>;
};
