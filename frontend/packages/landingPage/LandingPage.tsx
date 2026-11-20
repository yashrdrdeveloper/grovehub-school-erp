import { createContext, useContext, useEffect, useId, useMemo, useState, type FormEvent, type MouseEvent, type ReactNode } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import {
  ArrowRight, ArrowUpRight, BookOpen, Bus, CalendarDays,
  Check, Clock3, Computer, Download, ExternalLink, FileCheck2, FlaskConical,
  GraduationCap, HeartPulse, Mail, MapPin, Menu, Phone, Presentation,
  School, ShieldCheck, Sparkles, Trophy, Megaphone, CalendarClock, Pencil, Star, X,
  Landmark, HeartHandshake, Images,
} from 'lucide-react';

import { SchoolMonogram } from '@Components';

import { achievements, disclosedDocuments, enquiryTypes, facilities, historicDisclosure, photographs, school, upcomingEvents } from './content';
import { PublicSheet, type SchoolPanel } from './PublicSheet';

import './landingPage.css';

/** Explicit enquiry actions do not depend on catching React Router links
 * at the page root. Their href remains useful for opening a new tab or
 * loading the fully accessible admissions page without JavaScript. */
const EnquiryActionContext = createContext<(() => void) | null>(null);

const EnquiryLink = ({className,children,onClick}:{
  className?:string;children:ReactNode;onClick?:()=>void;
}) => {
  const openEnquiry=useContext(EnquiryActionContext);
  return <Link to="/admissions#enquiry" className={className} data-rdr-enquiry-action="true"
    onClick={event=>{
      onClick?.();
      if(event.defaultPrevented || event.button!==0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey)return;
      if(!openEnquiry)return;
      event.preventDefault();
      openEnquiry();
    }}>{children}</Link>;
};

const navItems = [
  { label:'About', to:'/about' },
  { label:'Academics', to:'/academics' },
  { label:'Campus', to:'/campus' },
  { label:'Admissions', to:'/admissions' },
  { label:'Gallery', to:'/gallery' },
  { label:'CBSE Disclosure', to:'/mandatory-public-disclosure' },
] as const;

const facilityIcons = [Computer, FlaskConical, Presentation, BookOpen, Trophy, Bus, HeartPulse, ShieldCheck];

type EnquiryFormData = {
  firstName:string;
  middleName:string;
  lastName:string;
  phone:string;
  email:string;
  type:string;
  interestedClass:string;
  message:string;
  privacy:boolean;
};

const emptyEnquiry: EnquiryFormData = {
  firstName:'', middleName:'', lastName:'', phone:'', email:'',
  type:'ADMISSION', interestedClass:'', message:'', privacy:false,
};

const ImageCard = ({
  src,alt,position,className='',loading='lazy',children,
}: {
  src:string;alt:string;position?:string;className?:string;loading?:'lazy'|'eager';children?:ReactNode;
}) => (
  <div className={`rdr-photo ${className}`}>
    <img src={src} alt={alt} loading={loading} decoding="async"
      style={{objectPosition:position ?? 'center'}} referrerPolicy="no-referrer"
      onError={event => {
        const image = event.currentTarget;
        const original = photographs.find(photo => photo.src === src)?.originalSrc;
        if (original && !image.dataset.fallbackAttempted) {
          image.dataset.fallbackAttempted = 'true';
          image.src = original;
          return;
        }
        // Avoid infinite onError loops and broken-image alt text on cards.
        image.style.display = 'none';
        image.parentElement?.classList.add('rdr-photo--unavailable');
      }} />
    {children}
  </div>
);

const SiteHeader = () => {
  const [open,setOpen] = useState(false);
  const { pathname } = useLocation();
  useEffect(() => setOpen(false),[pathname]);
  useEffect(() => {
    if (!open) return;
    const previous=document.body.style.overflow;
    document.body.style.overflow='hidden';
    return () => { document.body.style.overflow=previous; };
  },[open]);
  return (
    <>
      <div className="rdr-topline">
        <div className="rdr-container rdr-topline__inner">
          <span><MapPin size={13}/> Haluhera, Rewari, Haryana</span>
          <span><Phone size={13}/> <a href="tel:+919050077550">{school.phone}</a></span>
        </div>
      </div>
      <header className="rdr-header">
        <div className="rdr-container rdr-header__inner">
          <Link to="/" className="rdr-logo" aria-label="RDR School home">
            <SchoolMonogram/>
            <span><strong>RDR School</strong><small>Rao Deena Ram Vidya Vihar</small></span>
          </Link>
          <nav className="rdr-nav" aria-label="Main navigation">
            {navItems.map(item => <Link className={pathname===item.to?'active':''} to={item.to} key={item.to}>{item.label}</Link>)}
          </nav>
          <div className="rdr-header__actions">
            <EnquiryLink className="rdr-btn rdr-btn--solid rdr-header__enquire">Enquire now <ArrowUpRight size={16}/></EnquiryLink>
            <button type="button" className="rdr-menu-button" aria-label={open?'Close menu':'Open menu'}
              aria-expanded={open} aria-controls="rdr-mobile-navigation" onClick={()=>setOpen(v=>!v)}>
              {open?<X size={22}/>:<Menu size={22}/>}
            </button>
          </div>
        </div>
      </header>
      {open && <>
        <button className="rdr-mobile-nav__scrim" type="button" aria-label="Close menu" onClick={()=>setOpen(false)}/>
        <nav className="rdr-mobile-nav" id="rdr-mobile-navigation" aria-label="Mobile navigation">
          <div className="rdr-mobile-nav__top"><strong>Explore RDR School</strong><button type="button" onClick={()=>setOpen(false)} aria-label="Close menu"><X size={21}/></button></div>
          {navItems.map(item=><Link key={item.to} to={item.to} onClick={()=>setOpen(false)}>{item.label}<ArrowUpRight size={16}/></Link>)}
          <Link to="/about/director">Director's message <ArrowUpRight size={16}/></Link>
          <Link to="/about/principal">Principal's message <ArrowUpRight size={16}/></Link>
          <Link to="/contact">Contact & directions <ArrowUpRight size={16}/></Link>
          <EnquiryLink className="rdr-mobile-nav__action" onClick={()=>setOpen(false)}>Start an enquiry <ArrowRight size={16}/></EnquiryLink>
        </nav>
      </>}
    </>
  );
};

const SiteFooter = () => (
  <footer className="rdr-footer">
    <div className="rdr-container">
      <div className="rdr-footer__top">
        <div className="rdr-footer__intro">
          <div className="rdr-logo rdr-logo--footer"><SchoolMonogram/><span><strong>RDR School</strong><small>Rao Deena Ram Vidya Vihar</small></span></div>
          <p>Learning with purpose. Growing with confidence. Serving the community since 2006.</p>
          <span className="rdr-footer__location"><MapPin size={16}/>{school.address}</span>
        </div>
        <div><h3>Discover</h3><Link to="/about">Our story</Link><Link to="/academics">Academics</Link><Link to="/campus">Facilities</Link><Link to="/gallery">School life</Link></div>
        <div><h3>For families</h3><Link to="/admissions">Admissions</Link><Link to="/resources">Resources</Link><Link to="/contact">Contact</Link><Link to="/login">School workspace</Link></div>
        <div><h3>Transparency</h3><Link to="/mandatory-public-disclosure">Mandatory Public Disclosure</Link><a href={school.cbseSaras} target="_blank" rel="noreferrer">CBSE affiliation</a><a href={school.disclosurePdf} target="_blank" rel="noreferrer">2023 disclosure archive</a></div>
      </div>
      <div className="rdr-footer__bottom"><span>© {new Date().getFullYear()} RDR School. All rights reserved.</span><span>Haluhera · Rewari · Haryana</span></div>
    </div>
  </footer>
);

const Heading = ({eyebrow,title,description}:{eyebrow:string;title:string;description?:string}) => (
  <div className="rdr-section-heading"><span className="rdr-eyebrow">{eyebrow}</span><h2>{title}</h2>{description&&<p>{description}</p>}</div>
);

/** Purely decorative school motifs. Keep them outside important photo faces
 * and out of the accessibility tree; the original photographs remain intact. */
const PhotoMotifs = ({variant='hero'}:{variant?:'hero'|'gallery'|'banner'}) => (
  <div className={`rdr-photo-motifs rdr-photo-motifs--${variant}`} aria-hidden="true">
    <svg className="rdr-photo-motifs__flight" viewBox="0 0 170 76" fill="none">
      <path d="M4 59C34 26 48 83 83 44c20-21 25-38 51-36" stroke="currentColor" strokeWidth="1.6" strokeDasharray="4 7" strokeLinecap="round"/>
      <path d="m131 8 32-5-18 24-5-12-9-7Zm9 7 23-12" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round"/>
    </svg>
    <span className="rdr-photo-motifs__spark"><Sparkles size={27} strokeWidth={1.4}/></span>
    <span className="rdr-photo-motifs__star"><Star size={18} strokeWidth={1.5}/></span>
    <span className="rdr-photo-motifs__book"><BookOpen size={25} strokeWidth={1.4}/></span>
    <span className="rdr-photo-motifs__pencil"><Pencil size={23} strokeWidth={1.5}/></span>
  </div>
);

const PhotoCollage = () => (
  <div className="rdr-collage rdr-collage--editorial" aria-label="Authentic photographs from RDR School events">
    <ImageCard className="rdr-collage__large" src={photographs[0].src} alt={photographs[0].alt} position={photographs[0].position} loading="eager">
      <span className="rdr-collage__image-label">ACHIEVEMENTS <span>✳</span> RDR SCHOOL</span>
    </ImageCard>
    <ImageCard className="rdr-collage__top" src={photographs[1].src} alt={photographs[1].alt} position={photographs[1].position} loading="eager"/>
    <ImageCard className="rdr-collage__bottom" src={photographs[2].src} alt={photographs[2].alt} position={photographs[2].position} loading="eager"/>
    <PhotoMotifs/>
    <div className="rdr-collage__stamp"><Sparkles size={17}/> Every child, every possibility.</div>
  </div>
);

const EnquiryForm = ({compact=false}:{compact?:boolean}) => {
  const fieldId=useId();
  const [data,setData] = useState<EnquiryFormData>(emptyEnquiry);
  const [step,setStep] = useState(0);
  const [error,setError] = useState('');
  const [ready,setReady] = useState(false);
  const patch = (key:keyof EnquiryFormData,value:string|boolean) => setData(prev=>({...prev,[key]:value}));
  const next = (event:FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    if(!data.firstName.trim()){setError('Please enter the guardian’s first name.');return;}
    if(!data.phone.trim() && !data.email.trim()){setError('Enter a phone number or email address.');return;}
    if(data.email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email.trim())){setError('Enter a valid email address.');return;}
    if(data.phone.trim() && data.phone.replace(/\D/g,'').length<10){setError('Enter a valid phone number.');return;}
    setStep(1);
  };
  const submit = (event:FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    if(!data.privacy){setError('Please confirm you agree to be contacted about your enquiry.');return;}
    // No public Lead CRM endpoint exists in this branch. Avoid pretending that
    // a lead was created: provide a real email draft for the school instead.
    setReady(true);
  };
  const subject = encodeURIComponent(`RDR School — ${enquiryTypes.find(type=>type.value===data.type)?.label ?? 'Enquiry'}`);
  const body = encodeURIComponent([
    `Guardian: ${[data.firstName,data.middleName,data.lastName].filter(Boolean).join(' ')}`,
    `Phone: ${data.phone || 'Not provided'}`,
    `Email: ${data.email || 'Not provided'}`,
    `Enquiry: ${enquiryTypes.find(type=>type.value===data.type)?.label ?? data.type}`,
    `Interested class: ${data.interestedClass || 'Not provided'}`,
    `Message: ${data.message || '—'}`,
  ].join('\n'));
  return (
    <div className={`rdr-enquiry ${compact?'rdr-enquiry--compact':''}`} id={compact?undefined:'enquiry'}>
      <div className="rdr-enquiry__head"><span className="rdr-eyebrow">LET’S BEGIN</span><h3>{ready?'Your enquiry is ready': 'Have a question? We’re here.'}</h3><p>{ready?'Open your email application to send the enquiry to the school.': 'A quick conversation is the first step. We’ll help you find the answers.'}</p></div>
      {ready ? (
        <div className="rdr-enquiry__ready" role="status">
          <span className="rdr-enquiry__success"><Check size={22}/></span>
          <p>We’ve prepared your enquiry, but it has <strong>not been submitted</strong> to the school yet. Choose Send enquiry to open your email application.</p>
          <a className="rdr-btn rdr-btn--solid" href={`mailto:${school.email}?subject=${subject}&body=${body}`}>Send enquiry by email <Mail size={17}/></a>
          <button className="rdr-link-button" type="button" onClick={()=>{setReady(false);setStep(0);}}>Edit details</button>
        </div>
      ) : (
        <form onSubmit={step===0?next:submit} noValidate>
          <div className="rdr-enquiry__step"><span>STEP {step+1} OF 2</span><span>{step===0?'Your details':'A little more detail'}</span></div>
          {step===0 ? <>
            <label>Parent / guardian name <span>*</span></label>
            <div className="rdr-field-grid rdr-field-grid--three">
              <input aria-label="Guardian first name" placeholder="First name *" autoComplete="given-name" value={data.firstName} onChange={e=>patch('firstName',e.target.value)}/>
              <input aria-label="Guardian middle name" placeholder="Middle" value={data.middleName} onChange={e=>patch('middleName',e.target.value)}/>
              <input aria-label="Guardian last name" placeholder="Last name" autoComplete="family-name" value={data.lastName} onChange={e=>patch('lastName',e.target.value)}/>
            </div>
            <label htmlFor={`${fieldId}-phone`}>Phone number</label>
            <input id={`${fieldId}-phone`} type="tel" placeholder="+91 98765 43210" autoComplete="tel" value={data.phone} onChange={e=>patch('phone',e.target.value)}/>
            <label htmlFor={`${fieldId}-email`}>Email address</label>
            <input id={`${fieldId}-email`} type="email" placeholder="parent@example.com" autoComplete="email" value={data.email} onChange={e=>patch('email',e.target.value)}/>
            <small className="rdr-enquiry__hint">Please provide at least one: phone number or email.</small>
            <label htmlFor={`${fieldId}-type`}>I'm enquiring about <span>*</span></label>
            <select id={`${fieldId}-type`} value={data.type} onChange={e=>patch('type',e.target.value)}>
              {enquiryTypes.map(type=><option key={type.value} value={type.value}>{type.label}</option>)}
            </select>
          </>:<>
            {data.type==='ADMISSION' && <><label htmlFor={`${fieldId}-class`}>Class interested in</label><select id={`${fieldId}-class`} value={data.interestedClass} onChange={e=>patch('interestedClass',e.target.value)}>
              <option value="">Select a class (optional)</option>
              {['Nursery','KG','I','II','III','IV','V','VI','VII','VIII','IX','X','XI','XII'].map(value=><option key={value} value={value}>{value}</option>)}
            </select></>}
            <label htmlFor={`${fieldId}-message`}>Anything you'd like us to know?</label>
            <textarea id={`${fieldId}-message`} rows={4} placeholder="Tell us a little about your enquiry..." value={data.message} onChange={e=>patch('message',e.target.value)}/>
            <label className="rdr-enquiry__consent"><input type="checkbox" checked={data.privacy} onChange={e=>patch('privacy',e.target.checked)}/><span>I agree to share these details with RDR School to receive a response to my enquiry.</span></label>
          </>}
          {error && <p className="rdr-enquiry__error" role="alert">{error}</p>}
          <div className="rdr-enquiry__actions">
            {step===1 && <button className="rdr-btn rdr-btn--outline" type="button" onClick={()=>{setStep(0);setError('');}}>Back</button>}
            <button type="submit" className="rdr-btn rdr-btn--solid">{step===0?'Continue':'Prepare enquiry'} <ArrowRight size={17}/></button>
          </div>
        </form>
      )}
    </div>
  );
};

const StoriesAndUpdates = () => (
  <section className="rdr-section rdr-container rdr-stories">
    <div className="rdr-section__split">
      <Heading eyebrow="STAY CONNECTED" title="What's happening at RDR." description="Celebrating school life, memorable moments and the next opportunities for our community."/>
      <Link className="rdr-inline-link" to="/gallery">Explore our gallery <ArrowUpRight size={18}/></Link>
    </div>
    <div className="rdr-stories__grid">
      <article className="rdr-stories__featured">
        <ImageCard src={photographs[1].src} alt={photographs[1].alt} position={photographs[1].position}/>
        <div className="rdr-stories__featured-caption"><span>FROM OUR ARCHIVES · SCHOOL LIFE</span><h3>Learning beyond the classroom</h3><p>Celebrating expression, confidence and togetherness through school programmes.</p></div>
      </article>
      <div className="rdr-stories__side">
        <div className="rdr-stories__upcoming">
          <span className="rdr-eyebrow"><CalendarClock size={16}/> {upcomingEvents.length?'SCHOOL CALENDAR':'MEET RDR SCHOOL'}</span>
          <h3>{upcomingEvents.length?'Upcoming events':'The story behind our school'}</h3>
          {upcomingEvents.length ? upcomingEvents.map(event => (
            <div className="rdr-event" key={event.id}><strong>{event.date}</strong><span>{event.title}</span></div>
          )) : (
            <div className="rdr-school-message-links">
              <Link to="/about/history"><Landmark size={19}/><span><strong>Our school story</strong><small>Haluhera, since 2006</small></span><ArrowRight size={16}/></Link>
              <Link to="/about/director"><HeartHandshake size={19}/><span><strong>Director's message</strong><small>Care and all-round learning</small></span><ArrowRight size={16}/></Link>
              <Link to="/about/principal"><GraduationCap size={19}/><span><strong>Principal's message</strong><small>Education with character</small></span><ArrowRight size={16}/></Link>
            </div>
          )}
        </div>
        <EnquiryLink className="rdr-stories__admission">
          <span><Megaphone size={18}/> ADMISSIONS & VISITS</span>
          <strong>Ready to discover RDR School?</strong>
          <small>Ask about admissions, classes, transport and visiting the campus.</small>
          <ArrowUpRight size={21}/>
        </Link>
      </div>
    </div>
  </section>
);

const SchoolMessages = () => (
  <section className="rdr-section rdr-section--wash rdr-home-leaders">
    <div className="rdr-container">
      <Heading eyebrow="THE PEOPLE BEHIND OUR SCHOOL" title="A vision carried forward." description="Explore the school's history and messages from its leadership as published on the original RDR website."/>
      <div className="rdr-home-leaders__grid">
        <Link to="/about/history"><span>EST. 2006</span><Landmark size={28}/><h3>A school rooted in Haluhera.</h3><p>Discover where RDR began and the community it serves.</p><ArrowUpRight size={20}/></Link>
        <Link to="/about/director"><span>DIRECTOR'S MESSAGE</span><HeartHandshake size={28}/><h3>Education for every possibility.</h3><p>Shri Lal Singh Yadav on developing the whole child.</p><ArrowUpRight size={20}/></Link>
        <Link to="/about/principal"><span>PRINCIPAL'S MESSAGE</span><GraduationCap size={28}/><h3>Learning with character.</h3><p>Shri Arvind Yadav on education, discipline and growth.</p><ArrowUpRight size={20}/></Link>
      </div>
    </div>
  </section>
);

const MobileQuickHome = () => (
  <section className="rdr-quick-home rdr-container" aria-label="Quick school information and actions">
    <div className="rdr-quick-home__actions">
      <EnquiryLink className="rdr-quick-home__primary"><GraduationCap size={19}/><span>Admission enquiry</span><ArrowUpRight size={16}/></EnquiryLink>
      <a href="tel:+919050077550"><Phone size={19}/><span>Call school</span><ArrowUpRight size={16}/></a>
      <a href="https://www.google.com/maps/search/?api=1&query=Rao+Deena+Ram+Vidya+Vihar+Haluhera" target="_blank" rel="noreferrer"><MapPin size={19}/><span>Directions</span><ArrowUpRight size={16}/></a>
      <Link to="/mandatory-public-disclosure"><ShieldCheck size={19}/><span>CBSE details</span><ArrowUpRight size={16}/></Link>
    </div>
    <div className="rdr-quick-home__heading"><span className="rdr-eyebrow">DISCOVER OUR SCHOOL</span><h2>Everything important, a tap away.</h2></div>
    <div className="rdr-quick-home__tiles">
      {[
        {to:'/about/history',icon:Landmark,title:'Our story',desc:'Since 2006'},
        {to:'/about/director',icon:HeartHandshake,title:"Director's message",desc:'Growing every child'},
        {to:'/about/principal',icon:GraduationCap,title:"Principal's message",desc:'Learning with values'},
        {to:'/academics',icon:BookOpen,title:'Academics',desc:'Learning journey'},
        {to:'/campus',icon:School,title:'Campus',desc:'Spaces to explore'},
        {to:'/gallery',icon:Images,title:'School life',desc:'Moments in pictures'},
      ].map(item=>{
        const Icon=item.icon;
        return <Link to={item.to} key={item.to}><span className="rdr-quick-home__tile-icon"><Icon size={21}/></span><span><strong>{item.title}</strong><small>{item.desc}</small></span><ArrowUpRight size={15}/></Link>;
      })}
    </div>
    <div className="rdr-quick-home__bottom"><span><Sparkles size={16}/> A place to learn, belong and grow.</span><EnquiryLink>Talk to us <ArrowRight size={15}/></EnquiryLink></div>
  </section>
);

const HomePage = () => (
  <>
    <section className="rdr-hero rdr-container">
      <div className="rdr-hero__content">
        <span className="rdr-eyebrow"><span className="rdr-eyebrow__dot"/> EST. 2006 · HALUHERA, REWARI</span>
        <h1>Where curious minds <em>find their way.</em></h1>
        <p>Welcome to Rao Deena Ram Vidya Vihar — a school community that believes in thoughtful learning, strong values and the potential of every child.</p>
        <div className="rdr-hero__actions"><EnquiryLink className="rdr-btn rdr-btn--solid">Enquire for admissions <ArrowUpRight size={17}/></EnquiryLink><Link className="rdr-btn rdr-btn--text" to="/about">Get to know us <ArrowRight size={17}/></Link></div>
        <div className="rdr-hero__meta"><span><School size={17}/> CBSE affiliated</span><span><GraduationCap size={17}/> Senior Secondary</span><span><MapPin size={17}/> Rewari, Haryana</span></div>
      </div>
      <PhotoCollage/>
    </section>
    <MobileQuickHome/>
    <section className="rdr-ribbon"><div className="rdr-container rdr-ribbon__inner"><span>ROOTED IN VALUES</span><i>✳</i><span>GROWING THROUGH LEARNING</span><i>✳</i><span>TOGETHER AS A COMMUNITY</span></div></section>
    <section className="rdr-section rdr-container rdr-intro rdr-home-extended">
      <div className="rdr-intro__text"><Heading eyebrow="WELCOME TO RDR" title="A place to belong. A chance to become." description="Established in 2006 in Haluhera, our school has long served families from Rewari and surrounding villages. We aim to nurture students academically, morally and socially."/><Link to="/about" className="rdr-inline-link">Discover our story <ArrowUpRight size={18}/></Link></div>
      <div className="rdr-intro__photo"><ImageCard src={photographs[4].src} alt={photographs[4].alt}/><div className="rdr-intro__badge"><span>Since</span><strong>2006</strong><small>Committed to learning</small></div></div>
    </section>
    <section className="rdr-section rdr-section--wash rdr-home-facilities">
      <div className="rdr-container">
        <Heading eyebrow="EXPERIENCE RDR" title="Every day offers something new." description="The spaces, experiences and guidance that help young people learn, explore and grow."/>
        <div className="rdr-feature-grid">
          {facilities.slice(0,6).map((item,i)=>{
            const Icon=facilityIcons[i];
            return <div className="rdr-feature" key={item.title}><span className="rdr-feature__number">0{i+1}</span>{Icon&&<Icon size={27} strokeWidth={1.5}/>}<h3>{item.title}</h3><p>{item.description}</p></div>;
          })}
        </div>
        <Link className="rdr-inline-link" to="/campus">Explore the campus <ArrowUpRight size={18}/></Link>
      </div>
    </section>
    <section className="rdr-section rdr-container rdr-home-gallery">
      <div className="rdr-section__split"><Heading eyebrow="A SCHOOL FULL OF LIFE" title="The moments that make us." description="Real snapshots from the original RDR School photo gallery — celebrating creativity, achievements and friendship."/><Link to="/gallery" className="rdr-inline-link">See the gallery <ArrowUpRight size={18}/></Link></div>
      <div className="rdr-gallery-preview">
        {photographs.slice(1,5).map((photo,index)=><ImageCard key={photo.src} src={photo.src} alt={photo.alt} className={index===0?'rdr-gallery-preview__lead':''}>{<span>{photo.category}</span>}</ImageCard>)}
      </div>
    </section>
    <section className="rdr-achievements rdr-home-extended">
      <div className="rdr-container rdr-achievements__inner">
        {achievements.map(item => <div key={item.label}><strong>{item.value}</strong><span>{item.label}</span><small>{item.detail}</small></div>)}
      </div>
    </section>
    <div className="rdr-home-updates"><StoriesAndUpdates/></div>
    <div className="rdr-home-extended"><SchoolMessages/></div>
    <section className="rdr-section rdr-section--navy rdr-home-extended"><div className="rdr-container rdr-manifesto"><span className="rdr-eyebrow">OUR PURPOSE</span><h2>More than lessons.<br/><em>Learning for life.</em></h2><p>We believe education should help young people think critically, care for others and meet challenges with confidence.</p><Link className="rdr-btn rdr-btn--light" to="/academics">Explore academics <ArrowUpRight size={17}/></Link><span className="rdr-manifesto__ornament">✳</span></div></section>
    <section className="rdr-section rdr-container rdr-enquiry-section rdr-enquiry-section--short">
      <div><Heading eyebrow="ADMISSIONS & ENQUIRIES" title="A new chapter starts with a conversation." description="Ask about admissions, classes, transport or visiting RDR. Our team is here to help."/></div>
      <EnquiryLink className="rdr-btn rdr-btn--solid">Open quick enquiry <ArrowUpRight size={17}/></EnquiryLink>
    </section>
    <section className="rdr-section rdr-container rdr-disclosure-teaser"><div><span className="rdr-eyebrow">SCHOOL TRANSPARENCY</span><h2>Mandatory Public Disclosure</h2><p>Explore the school's CBSE affiliation reference, archived documents, and disclosure information.</p></div><Link to="/mandatory-public-disclosure" className="rdr-btn rdr-btn--outline">View disclosure <FileCheck2 size={18}/></Link></section>
  </>
);

const IntroPage = ({eyebrow,title,description,image}:{eyebrow:string;title:string;description:string;image?:string}) => (
  <section className="rdr-page-banner rdr-container"><div><span className="rdr-eyebrow">{eyebrow}</span><h1>{title}</h1><p>{description}</p></div>{image&&<div className="rdr-page-banner__art"><ImageCard className="rdr-page-banner__image" src={image} alt="Photograph of students and school activities at RDR School"/><PhotoMotifs variant="banner"/></div>}</section>
);

const AboutPage = () => (
  <>
    <IntroPage eyebrow="OUR STORY" title="A school with purpose, rooted in community." description="Rao Deena Ram Vidya Vihar was established in Haluhera in 2006, dedicated to thoughtful education and the all-round development of students." image={photographs[4].src}/>
    <section className="rdr-section rdr-container rdr-prose-grid">
      <div><span className="rdr-eyebrow">OUR BEGINNINGS</span><h2>Inspired by a legacy.</h2></div>
      <div><p>The school is named in memory of Late Shri Rao Deena Ram Ji. According to the school's history, it was founded in April 2006 by Shri Lal Singh Yadav (Retired HES-II), to extend educational opportunities to the surrounding rural community.</p><p>Its stated philosophy emphasizes respect, responsibility, fairness, caring and citizenship. These values continue to shape how the school's story is presented.</p><a href="https://rdrschool.com/our-school.html" target="_blank" rel="noreferrer" className="rdr-inline-link">Read the original school history <ExternalLink size={16}/></a></div>
    </section>
    <section className="rdr-section rdr-section--wash"><div className="rdr-container rdr-values">
      {[{number:'01',title:'Learning',description:'Encouraging curiosity, creativity and independent thinking.'},{number:'02',title:'Character',description:'Building a foundation of responsibility, kindness and integrity.'},{number:'03',title:'Community',description:'Strengthening partnerships between families, teachers and students.'}].map(x=><div key={x.number}><span>{x.number}</span><h3>{x.title}</h3><p>{x.description}</p></div>)}
    </div></section>
    <section className="rdr-section rdr-container rdr-about-leaders">
      <Heading eyebrow="WORDS THAT GUIDE US" title="Meet the people and values behind RDR."/>
      <div className="rdr-about-leaders__links">
        <Link to="/about/history"><Landmark size={22}/><strong>Our history</strong><small>From Haluhera, since 2006</small><ArrowUpRight size={16}/></Link>
        <Link to="/about/director"><HeartHandshake size={22}/><strong>Director's message</strong><small>All-round development</small><ArrowUpRight size={16}/></Link>
        <Link to="/about/principal"><GraduationCap size={22}/><strong>Principal's message</strong><small>Learning with character</small><ArrowUpRight size={16}/></Link>
        <Link to="/about/philosophy"><Sparkles size={22}/><strong>Our philosophy</strong><small>Broadening minds</small><ArrowUpRight size={16}/></Link>
      </div>
    </section>
    <section className="rdr-section rdr-container">
      <Heading eyebrow="LEADERSHIP & PHILOSOPHY" title="Guided by people who care." description="The school's historical leadership messages speak about holistic learning, kindness and preparing children for life beyond the classroom."/>
      <div className="rdr-values rdr-values--leadership">
        <div><span>FOUNDER & DIRECTOR · HISTORICAL PROFILE</span><h3>Shri Lal Singh Yadav</h3><p>The original site credits Shri Lal Singh Yadav (Retired HES-II) with founding the school in April 2006. His published message emphasizes students' all-round development.</p><a className="rdr-inline-link" href="https://rdrschool.com/director-message.html" target="_blank" rel="noreferrer">Original director's message <ExternalLink size={16}/></a></div>
        <div><span>PRINCIPAL · 2023 DISCLOSURE</span><h3>Arvind Yadav</h3><p>Named as principal in the school's 2023 CBSE Appendix IX. Current office holder and leadership portrait must be verified before publication.</p><a className="rdr-inline-link" href="https://rdrschool.com/Principal-message.html" target="_blank" rel="noreferrer">Original principal's message <ExternalLink size={16}/></a></div>
        <div><span>EDUCATIONAL PHILOSOPHY</span><h3>Be a lamp unto yourself</h3><p>The legacy school describes its emblem as a lamp symbolizing personal understanding and responsibility. We retain that original identity alongside the RDR digital wordmark.</p><a className="rdr-inline-link" href="https://rdrschool.com/our-philosophy.html" target="_blank" rel="noreferrer">School philosophy <ExternalLink size={16}/></a></div>
      </div>
    </section>
  </>
);

/** Direct-addressable, crawlable pages for school leadership and history.
 * In-app navigation opens the visual panel instead of loading these pages. */
const legacySchoolStories = {
  history:{
    eyebrow:'ESTABLISHED · APRIL 2006',
    title:'Our school story.',
    lead:'A school rooted in the community of Haluhera and built around the promise of better learning opportunities.',
    detail:'Rao Deena Ram Vidya Vihar was named in memory of Late Shri Rao Deena Ram Ji. According to the original school history, Shri Lal Singh Yadav (Retired HES-II) founded the school in April 2006. Its identity reflects a continuing relationship between families, teachers and the local community.',
    highlights:['A school founded in 2006','Haluhera, District Rewari','Education and community partnership'],
    src:'https://rdrschool.com/our-school.html',
    image:4,
  },
  director:{
    eyebrow:'A MESSAGE FROM OUR DIRECTOR · LEGACY SCHOOL SITE',
    title:"Director's message.",
    lead:'Every child deserves opportunities to learn, express themselves and grow as a person.',
    detail:"In his original school message, Shri Lal Singh Yadav describes a vision of all-round development. The school's commitment includes academics, sports, arts, drama and co-curricular activities, alongside strong values and preparation for each child's future.",
    highlights:['All-round development','Values and education','Arts, sports and academics'],
    src:'https://rdrschool.com/director-message.html',
    image:0,
  },
  principal:{
    eyebrow:'A MESSAGE FROM OUR PRINCIPAL · LEGACY SCHOOL SITE',
    title:"Principal's message.",
    lead:'Confidence, discipline and character are as important as academic progress.',
    detail:"The principal's historical message, attributed to Shri Arvind Yadav, recounts the school's beginnings and growth. It highlights laboratories, experienced teachers, co-curricular opportunities and character-building education. The current office-holder should be reconfirmed before official publication.",
    highlights:['Learning with discipline','Experienced teaching','Character and responsibility'],
    src:'https://rdrschool.com/Principal-message.html',
    image:0,
  },
  philosophy:{
    eyebrow:'OUR VISION · BROADENING MINDS',
    title:'Our philosophy.',
    lead:'Be a lamp unto yourself.',
    detail:"RDR School's original philosophy emphasizes a creative and caring environment. It seeks to nurture students intellectually, physically, socially and emotionally, encouraging independence, compassion, honesty and critical thinking alongside academic learning.",
    highlights:['Independent thinking','Care and compassion','Learning for life'],
    src:'https://rdrschool.com/our-philosophy.html',
    image:2,
  },
} as const;

const StoryPage = ({topic}:{topic:keyof typeof legacySchoolStories}) => {
  const story=legacySchoolStories[topic];
  const image=photographs[story.image];
  return <>
    <IntroPage eyebrow={story.eyebrow} title={story.title} description={story.lead} image={image?.src}/>
    <section className="rdr-section rdr-container rdr-story-page">
      <div><span className="rdr-eyebrow">A CLOSER LOOK</span><h2>Inside the RDR community.</h2><p>{story.detail}</p>
        <div className="rdr-story-page__actions"><EnquiryLink className="rdr-btn rdr-btn--solid">Ask the school <ArrowUpRight size={16}/></EnquiryLink><a href={story.src} target="_blank" rel="noreferrer" className="rdr-inline-link">Original school message <ExternalLink size={16}/></a></div>
      </div>
      <div className="rdr-story-page__notes">{story.highlights.map((line,i)=><div key={line}><span>0{i+1}</span><strong>{line}</strong><Check size={16}/></div>)}</div>
    </section>
  </>;
};

const AcademicsPage = () => (
  <>
    <IntroPage eyebrow="ACADEMICS" title="Learning that shapes a lifetime." description="A broad educational experience focused on understanding, confidence, creativity and practical learning." image={photographs[0].src}/>
    <section className="rdr-section rdr-container"><Heading eyebrow="LEARNING JOURNEY" title="Room to think. Space to discover." description="From the earliest years through senior secondary school, our focus is on building strong foundations and supporting each learner's growth."/>
      <div className="rdr-feature-grid">{[
        ['Foundational Years','Exploration, play and early learning experiences.'],
        ['Primary Years','Building literacy, numeracy and confidence.'],
        ['Middle School','Developing independence and deeper understanding.'],
        ['Secondary & Senior Secondary','Preparing students for advanced study and future pathways.'],
      ].map(([title,description],i)=><div key={title} className="rdr-feature"><span className="rdr-feature__number">0{i+1}</span><BookOpen size={28}/><h3>{title}</h3><p>{description}</p></div>)}</div>
    </section>
    <section className="rdr-section rdr-section--wash"><div className="rdr-container rdr-prose-grid"><div><Heading eyebrow="BEYOND THE CLASSROOM" title="Learning in action."/></div><p>Our historical school information describes computer education, science facilities, reading resources, sports, cultural programmes and student activities. For current class offerings, streams, subject combinations and curricula, please contact the school office.</p></div></section>
  </>
);

const CampusPage = () => (
  <>
    <IntroPage eyebrow="CAMPUS & FACILITIES" title="Spaces for discovery and belonging." description="From learning laboratories to playful outdoor moments, explore the facilities described on the school's original website." image={photographs[2].src}/>
    <section className="rdr-section rdr-container"><Heading eyebrow="FACILITIES" title="Designed for a fuller school day."/>
      <div className="rdr-feature-grid">{facilities.map((x,i)=>{const Icon=facilityIcons[i];return <div className="rdr-feature" key={x.title}>{Icon&&<Icon size={29} strokeWidth={1.5}/>}<h3>{x.title}</h3><p>{x.description}</p></div>;})}</div>
      <div className="rdr-campus-spotlight"><ImageCard src={photographs[2].src} alt={photographs[2].alt}/><div><span className="rdr-eyebrow">CAMPUS LIFE</span><h3>Spaces to explore, memories to make.</h3><p>From early experiences to senior school, each day brings moments of curiosity, teamwork and discovery.</p><Link className="rdr-inline-link" to="/gallery">See school life <ArrowUpRight size={17}/></Link></div></div>
    </section>
  </>
);

const GalleryPage = () => (
  <>
    <IntroPage eyebrow="GALLERY" title="Life at RDR, in moments." description="A photo story of school celebrations, student life and achievements — captured at RDR School."/>
    <section className="rdr-section rdr-container">
      <div className="rdr-gallery-grid">{photographs.map((photo,i)=><figure className={i%5===0?'rdr-gallery-grid__wide':''} key={photo.src}><div className="rdr-gallery-grid__art"><ImageCard src={photo.src} alt={photo.alt} position={photo.position}/>{i%3===0&&<PhotoMotifs variant="gallery"/>}</div><figcaption><span>{photo.category}</span><strong>{photo.title}</strong></figcaption></figure>)}</div>
      <p className="rdr-gallery__archive-note">From the RDR School photo archives. <a href="https://rdrschool.com/gallery.html" target="_blank" rel="noreferrer">See the original gallery <ExternalLink size={13}/></a></p>
    </section>
  </>
);

const AdmissionsPage = () => (
  <>
    <IntroPage eyebrow="JOIN OUR COMMUNITY" title="Your child's journey starts here." description="Exploring school options? We can help with admissions, school visits, fees, transport and academic questions." image={photographs[3].src}/>
    <section className="rdr-section rdr-container rdr-admissions-layout">
      <div><Heading eyebrow="HOW TO BEGIN" title="An easier first step." description="The school's historical admission procedure asks families to contact the school office and collect an application. Current eligibility, seat availability and the academic-year schedule should be verified directly."/>
        <div className="rdr-steps">{[['01','Share an enquiry','Tell us what you would like to know.'],['02','Talk to the school','Discuss class availability, visits and requirements.'],['03','Follow the admission process','Get current forms and guidance directly from the school.']].map(([n,t,d])=><div key={n}><span>{n}</span><div><h3>{t}</h3><p>{d}</p></div></div>)}</div>
        <div className="rdr-admissions__highlight"><Sparkles size={20}/><span>Every great school journey starts with a conversation. Our team is here to guide you through the next step.</span></div>
      </div>
      <EnquiryForm/>
    </section>
  </>
);

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

const pageContent: Record<string,()=>ReactNode> = {
  '/':HomePage,
  '/about':AboutPage,
  '/about/history':()=> <StoryPage topic="history"/>,
  '/about/director':()=> <StoryPage topic="director"/>,
  '/about/principal':()=> <StoryPage topic="principal"/>,
  '/about/philosophy':()=> <StoryPage topic="philosophy"/>,
  '/academics':AcademicsPage,
  '/campus':CampusPage,
  '/admissions':AdmissionsPage,
  '/gallery':GalleryPage,
  '/mandatory-public-disclosure':DisclosurePage,
  '/resources':ResourcesPage,
  '/contact':ContactPage,
};

const pathPanel: Record<string,SchoolPanel> = {
  '/about':'about',
  '/about/history':'history',
  '/about/director':'director',
  '/about/principal':'principal',
  '/about/philosophy':'philosophy',
  '/academics':'academics',
  '/campus':'campus',
  '/admissions':'enquiry',
  '/gallery':'gallery',
  '/resources':'resources',
  '/contact':'contact',
};

const pageMetadata: Record<string,{title:string;description:string}> = {
  '/':{title:'RDR School | Rao Deena Ram Vidya Vihar, Haluhera',description:'Discover Rao Deena Ram Vidya Vihar School in Haluhera, Rewari. Learn about school life, academic programmes, facilities and admissions.'},
  '/about':{title:'About RDR School | History, Vision & Leadership',description:'Explore RDR School history since 2006, its philosophy, director and principal messages in Haluhera, Rewari.'},
  '/about/history':{title:'School History | RDR School Haluhera',description:'Discover how RDR School was founded in April 2006 in Haluhera, Rewari, by Shri Lal Singh Yadav.'},
  '/about/director':{title:"Director's Message | RDR School",description:"Discover Shri Lal Singh Yadav's published message about all-round student development at RDR School."},
  '/about/principal':{title:"Principal's Message | RDR School",description:"Read the historical RDR School principal's message about academic learning, character and school values."},
  '/about/philosophy':{title:'School Vision & Philosophy | RDR School',description:"Learn about RDR School's vision of broadening minds, independent thinking, compassion and growth."},
  '/academics':{title:'Academics | RDR School Haluhera',description:'Explore learning stages, academics and school activities at Rao Deena Ram Vidya Vihar in Rewari.'},
  '/campus':{title:'Campus & Facilities | RDR School',description:'Explore the school campus, science labs, computer lab, smart classrooms, sports and student facilities.'},
  '/admissions':{title:'Admissions & Enquiry | RDR School',description:'Contact RDR School about admissions, classes, fees, transport and visits. Submit a quick enquiry to the school office.'},
  '/gallery':{title:'School Life & Gallery | RDR School',description:'See moments of student life, culture, sports and school achievements from the RDR School gallery.'},
  '/mandatory-public-disclosure':{title:'CBSE Mandatory Public Disclosure | RDR School',description:'Find RDR School CBSE affiliation details, historic Appendix IX disclosures and school documents.'},
  '/resources':{title:'School Resources | RDR School',description:'Find school academic information, the calendar and available documents from RDR School.'},
  '/contact':{title:'Contact RDR School | Haluhera, Rewari',description:'Contact Rao Deena Ram Vidya Vihar School in Haluhera, Rewari, Haryana. Phone, email, directions and enquiries.'},
};

export const LandingPage = () => {
  const location=useLocation();
  const navigate=useNavigate();
  const [enquiryOpen,setEnquiryOpen]=useState(false);
  const state=location.state as {background?:string;panel?:SchoolPanel}|null;
  const routedPanel=state?.background && state.panel ? state.panel : null;
  const panel:SchoolPanel|null=enquiryOpen?'enquiry':routedPanel;
  const backgroundPath=routedPanel?(state?.background??'/'):location.pathname;
  const CurrentPage=useMemo(()=>pageContent[backgroundPath] ?? HomePage,[backgroundPath]);

  const closePanel=()=>{
    if(enquiryOpen){setEnquiryOpen(false);return;}
    navigate(state?.background??'/',{replace:true});
  };
  const selectPanel=(next:SchoolPanel)=>{
    if(next==='enquiry'){setEnquiryOpen(true);return;}
    setEnquiryOpen(false);
    const path=Object.entries(pathPanel).find(([,value])=>value===next)?.[0]??'/about';
    navigate(path,{replace:true,state:{background:backgroundPath,panel:next}});
  };

  // A real URL is always retained for navigation, direct links and SEO.
  // Clicking around the public site opens a sheet without leaving the
  // visible landing page; loading a URL directly renders its full HTML page.
  const handleInternalLink=(event:MouseEvent<HTMLDivElement>)=>{
    if(event.defaultPrevented || event.button!==0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey)return;
    if(!(event.target instanceof Element))return;
    const anchor=event.target.closest<HTMLAnchorElement>('a[href]');
    if(!anchor || anchor.hasAttribute('download') || (anchor.target&&anchor.target!=='_self'))return;
    // Purpose-built enquiry links handle themselves. For generic Admissions
    // links (header menu/footer) preserve the same working enquiry behavior.
    if(anchor.dataset.rdrEnquiryAction==='true')return;
    const url=new URL(anchor.href,window.location.href);
    if(url.origin!==window.location.origin)return;
    if(url.pathname==='/admissions'){
      event.preventDefault();
      setEnquiryOpen(true);
      return;
    }
    const requested=pathPanel[url.pathname];
    if(!requested)return; // CBSE disclosures stay on their public HTML URL.
    event.preventDefault();
    navigate(url.pathname+url.search+url.hash,{
      state:{background:backgroundPath,panel:requested},
    });
  };
  useEffect(()=>{
    if(panel)return;
    if(location.hash){
      window.requestAnimationFrame(()=>document.getElementById(location.hash.slice(1))?.scrollIntoView({behavior:'smooth',block:'start'}));
    } else window.scrollTo({top:0,behavior:'instant'});
  },[backgroundPath,location.hash,panel]);
  useEffect(()=>{
    const details=pageMetadata[location.pathname] ?? (
      location.pathname.startsWith('/about/')?pageMetadata['/about']:pageMetadata['/']
    );
    if(!details)return;
    document.title=details.title;
    const description=document.querySelector<HTMLMetaElement>('meta[name="description"]');
    if(description)description.content=details.description;
    let canonical=document.querySelector<HTMLLinkElement>('link[rel="canonical"]');
    if(!canonical){canonical=document.createElement('link');canonical.rel='canonical';document.head.appendChild(canonical);}
    canonical.href=`https://rdrschool.com${location.pathname}`;
  },[location.pathname]);

  return <EnquiryActionContext.Provider value={()=>setEnquiryOpen(true)}>
    <div className="rdr-site" onClickCapture={handleInternalLink}>
      <SiteHeader/>
      <main id="main-content"><CurrentPage/></main>
      <SiteFooter/>
      <PublicSheet panel={panel} onClose={closePanel} onNavigate={selectPanel} enquiry={<EnquiryForm compact/>}/>
    </div>
  </EnquiryActionContext.Provider>;
};
