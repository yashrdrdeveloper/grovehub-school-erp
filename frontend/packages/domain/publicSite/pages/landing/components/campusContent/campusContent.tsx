import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { Link } from 'react-router-dom';
import { ArrowUpRight, BookOpen, ExternalLink, Sparkles } from '@Components/icons';
import { facilities, photographs } from '../../constants/school.content';
import { Heading, ImageCard, PhotoMotifs, facilityIcons } from '../siteMedia';
import { IntroPage } from '../aboutContent';
import { EnquiryForm } from '../enquiryForm';
import { getGalleryPhotoClassName } from './helpers/campusContent.getGalleryPhotoClassName';

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
      <div className="rdr-campus-spotlight"><ImageCard src={photographs[2].src} alt={photographs[2].alt}/><div><span className="rdr-eyebrow">CAMPUS LIFE</span><h3>Spaces to explore, memories to make.</h3><p>From early experiences to senior school, each day brings moments of curiosity, teamwork and discovery.</p><Link className="rdr-inline-link" to={SCHOOL_ROUTE.GALLERY}>See school life <ArrowUpRight size={17}/></Link></div></div>
    </section>
  </>
);

const GalleryPage = () => (
  <>
    <IntroPage eyebrow="GALLERY" title="Life at RDR, in moments." description="A photo story of school celebrations, student life and achievements — captured at RDR School."/>
    <section className="rdr-section rdr-container">
      <div className="rdr-gallery-grid">{photographs.map((photo,i)=><figure className={getGalleryPhotoClassName(i)} key={photo.src}><div className="rdr-gallery-grid__art"><ImageCard src={photo.src} alt={photo.alt} position={photo.position}/>{i%3===0&&<PhotoMotifs variant="gallery"/>}</div><figcaption><span>{photo.category}</span><strong>{photo.title}</strong></figcaption></figure>)}</div>
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

export { AcademicsPage };
export { CampusPage };
export { GalleryPage };
export { AdmissionsPage };
