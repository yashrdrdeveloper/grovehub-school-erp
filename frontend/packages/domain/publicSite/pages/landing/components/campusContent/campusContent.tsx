import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Box } from '@Components/components/Box';
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
    <Box as="section" className="rdr-section rdr-container"><Heading eyebrow="LEARNING JOURNEY" title="Room to think. Space to discover." description="From the earliest years through senior secondary school, our focus is on building strong foundations and supporting each learner's growth."/>
      <Box as="div" className="rdr-feature-grid">{[
        ['Foundational Years','Exploration, play and early learning experiences.'],
        ['Primary Years','Building literacy, numeracy and confidence.'],
        ['Middle School','Developing independence and deeper understanding.'],
        ['Secondary & Senior Secondary','Preparing students for advanced study and future pathways.'],
      ].map(([title,description],i)=><Box as="div" key={title} className="rdr-feature"><Box as="span" className="rdr-feature__number">0{i+1}</Box><BookOpen size={28}/><Box as="h3">{title}</Box><Box as="p">{description}</Box></Box>)}</Box>
    </Box>
    <Box as="section" className="rdr-section rdr-section--wash"><Box as="div" className="rdr-container rdr-prose-grid"><Box as="div"><Heading eyebrow="BEYOND THE CLASSROOM" title="Learning in action."/></Box><Box as="p">Our historical school information describes computer education, science facilities, reading resources, sports, cultural programmes and student activities. For current class offerings, streams, subject combinations and curricula, please contact the school office.</Box></Box></Box>
  </>
);

const CampusPage = () => (
  <>
    <IntroPage eyebrow="CAMPUS & FACILITIES" title="Spaces for discovery and belonging." description="From learning laboratories to playful outdoor moments, explore the facilities described on the school's original website." image={photographs[2].src}/>
    <Box as="section" className="rdr-section rdr-container"><Heading eyebrow="FACILITIES" title="Designed for a fuller school day."/>
      <Box as="div" className="rdr-feature-grid">{facilities.map((x,i)=>{const Icon=facilityIcons[i];return <Box as="div" className="rdr-feature" key={x.title}>{Icon&&<Icon size={29} strokeWidth={1.5}/>}<Box as="h3">{x.title}</Box><Box as="p">{x.description}</Box></Box>;})}</Box>
      <Box as="div" className="rdr-campus-spotlight"><ImageCard src={photographs[2].src} alt={photographs[2].alt}/><Box as="div"><Box as="span" className="rdr-eyebrow">CAMPUS LIFE</Box><Box as="h3">Spaces to explore, memories to make.</Box><Box as="p">From early experiences to senior school, each day brings moments of curiosity, teamwork and discovery.</Box><Link className="rdr-inline-link" to={SCHOOL_ROUTE.GALLERY}>See school life <ArrowUpRight size={17}/></Link></Box></Box>
    </Box>
  </>
);

const GalleryPage = () => (
  <>
    <IntroPage eyebrow="GALLERY" title="Life at RDR, in moments." description="A photo story of school celebrations, student life and achievements — captured at RDR School."/>
    <Box as="section" className="rdr-section rdr-container">
      <Box as="div" className="rdr-gallery-grid">{photographs.map((photo,i)=><figure className={getGalleryPhotoClassName(i)} key={photo.src}><Box as="div" className="rdr-gallery-grid__art"><ImageCard src={photo.src} alt={photo.alt} position={photo.position}/>{i%3===0&&<PhotoMotifs variant="gallery"/>}</Box><figcaption><Box as="span">{photo.category}</Box><Box as="strong">{photo.title}</Box></figcaption></figure>)}</Box>
      <Box as="p" className="rdr-gallery__archive-note">From the RDR School photo archives. <ExternalAnchor href="https://rdrschool.com/gallery.html" target="_blank" rel="noreferrer">See the original gallery <ExternalLink size={13}/></ExternalAnchor></Box>
    </Box>
  </>
);

const AdmissionsPage = () => (
  <>
    <IntroPage eyebrow="JOIN OUR COMMUNITY" title="Your child's journey starts here." description="Exploring school options? We can help with admissions, school visits, fees, transport and academic questions." image={photographs[3].src}/>
    <Box as="section" className="rdr-section rdr-container rdr-admissions-layout">
      <Box as="div"><Heading eyebrow="HOW TO BEGIN" title="An easier first step." description="The school's historical admission procedure asks families to contact the school office and collect an application. Current eligibility, seat availability and the academic-year schedule should be verified directly."/>
        <Box as="div" className="rdr-steps">{[['01','Share an enquiry','Tell us what you would like to know.'],['02','Talk to the school','Discuss class availability, visits and requirements.'],['03','Follow the admission process','Get current forms and guidance directly from the school.']].map(([n,t,d])=><Box as="div" key={n}><Box as="span">{n}</Box><Box as="div"><Box as="h3">{t}</Box><Box as="p">{d}</Box></Box></Box>)}</Box>
        <Box as="div" className="rdr-admissions__highlight"><Sparkles size={20}/><Box as="span">Every great school journey starts with a conversation. Our team is here to guide you through the next step.</Box></Box>
      </Box>
      <EnquiryForm/>
    </Box>
  </>
);

export { AcademicsPage };
export { CampusPage };
export { GalleryPage };
export { AdmissionsPage };
