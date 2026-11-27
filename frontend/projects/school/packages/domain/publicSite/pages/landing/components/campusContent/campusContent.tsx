import { _map } from '@Base/utilities/map';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Box , BOX_TAG} from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { AppLink } from '@Components/components/AppLink';
import { ArrowUpRight, BookOpen, ExternalLink, Sparkles } from '@Components/icons';
import { facilities, photographs } from '../../constants/school.content';
import { Heading, ImageCard, PhotoMotifs, PHOTO_MOTIFS_VARIANT, facilityIcons } from '../siteMedia';
import { IntroPage } from '../aboutContent';
import { EnquiryForm } from '../enquiryForm';
import { getGalleryPhotoClassName } from './helpers/campusContent.getGalleryPhotoClassName';

const AcademicsPage = () => (
  <>
    <IntroPage eyebrow="ACADEMICS" title={__("Learning that shapes a lifetime.")} description={__("A broad educational experience focused on understanding, confidence, creativity and practical learning.")} image={photographs[0].src}/>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container"><Heading eyebrow="LEARNING JOURNEY" title={__("Room to think. Space to discover.")} description={__("From the earliest years through senior secondary school, our focus is on building strong foundations and supporting each learner's growth.")}/>
      <Box as={BOX_TAG.DIV} className="rdr-feature-grid">{_map([
        ['Foundational Years','Exploration, play and early learning experiences.'],
        ['Primary Years','Building literacy, numeracy and confidence.'],
        ['Middle School','Developing independence and deeper understanding.'],
        ['Secondary & Senior Secondary','Preparing students for advanced study and future pathways.'],
      ], ([title,description],i)=><Box as={BOX_TAG.DIV} key={title} className="rdr-feature"><Box as={BOX_TAG.SPAN} className="rdr-feature__number">0{i+1}</Box><BookOpen size={28}/><Box as={BOX_TAG.H3}>{title}</Box><Box as={BOX_TAG.P}>{description}</Box></Box>)}</Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-section--wash"><Box as={BOX_TAG.DIV} className="rdr-container rdr-prose-grid"><Box as={BOX_TAG.DIV}><Heading eyebrow="BEYOND THE CLASSROOM" title={__("Learning in action.")}/></Box><Box as={BOX_TAG.P}>{__("Our historical school information describes computer education, science facilities, reading resources, sports, cultural programmes and student activities. For current class offerings, streams, subject combinations and curricula, please contact the school office.")}</Box></Box></Box>
  </>
);

const CampusPage = () => (
  <>
    <IntroPage eyebrow="CAMPUS & FACILITIES" title={__("Spaces for discovery and belonging.")} description={__("From learning laboratories to playful outdoor moments, explore the facilities described on the school's original website.")} image={photographs[2].src}/>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container"><Heading eyebrow="FACILITIES" title={__("Designed for a fuller school day.")}/>
      <Box as={BOX_TAG.DIV} className="rdr-feature-grid">{_map(facilities, (x,i)=>{const Icon=facilityIcons[i];return <Box as={BOX_TAG.DIV} className="rdr-feature" key={x.title}>{Icon&&<Icon size={29} strokeWidth={1.5}/>}<Box as={BOX_TAG.H3}>{x.title}</Box><Box as={BOX_TAG.P}>{x.description}</Box></Box>;})}</Box>
      <Box as={BOX_TAG.DIV} className="rdr-campus-spotlight"><ImageCard src={photographs[2].src} alt={photographs[2].alt}/><Box as={BOX_TAG.DIV}><Box as={BOX_TAG.SPAN} className="rdr-eyebrow">{__("CAMPUS LIFE")}</Box><Box as={BOX_TAG.H3}>{__("Spaces to explore, memories to make.")}</Box><Box as={BOX_TAG.P}>{__("From early experiences to senior school, each day brings moments of curiosity, teamwork and discovery.")}</Box><AppLink className="rdr-inline-link" to={SCHOOL_ROUTE.GALLERY}>{__("See school life ")}<ArrowUpRight size={17}/></AppLink></Box></Box>
    </Box>
  </>
);

const GalleryPage = () => (
  <>
    <IntroPage eyebrow="GALLERY" title={__("Life at RDR, in moments.")} description={__("A photo story of school celebrations, student life and achievements — captured at RDR School.")}/>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container">
      <Box as={BOX_TAG.DIV} className="rdr-gallery-grid">{_map(photographs, (photo,i)=><figure className={getGalleryPhotoClassName(i)} key={photo.src}><Box as={BOX_TAG.DIV} className="rdr-gallery-grid__art"><ImageCard src={photo.src} alt={photo.alt} position={photo.position}/>{i%3===0&&<PhotoMotifs variant={PHOTO_MOTIFS_VARIANT.GALLERY}/>}</Box><figcaption><Box as={BOX_TAG.SPAN}>{photo.category}</Box><Box as={BOX_TAG.STRONG}>{photo.title}</Box></figcaption></figure>)}</Box>
      <Box as={BOX_TAG.P} className="rdr-gallery__archive-note">{__("From the RDR School photo archives. ")}<ExternalAnchor href="https://rdrschool.com/gallery.html" target="_blank" rel="noreferrer">{__("See the original gallery ")}<ExternalLink size={13}/></ExternalAnchor></Box>
    </Box>
  </>
);

const AdmissionsPage = () => (
  <>
    <IntroPage eyebrow="JOIN OUR COMMUNITY" title={__("Your child's journey starts here.")} description={__("Exploring school options? We can help with admissions, school visits, fees, transport and academic questions.")} image={photographs[3].src}/>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container rdr-admissions-layout">
      <Box as={BOX_TAG.DIV}><Heading eyebrow="HOW TO BEGIN" title={__("An easier first step.")} description={__("The school's historical admission procedure asks families to contact the school office and collect an application. Current eligibility, seat availability and the academic-year schedule should be verified directly.")}/>
        <Box as={BOX_TAG.DIV} className="rdr-steps">{_map([['01','Share an enquiry','Tell us what you would like to know.'],['02','Talk to the school','Discuss class availability, visits and requirements.'],['03','Follow the admission process','Get current forms and guidance directly from the school.']], ([n,t,d])=><Box as={BOX_TAG.DIV} key={n}><Box as={BOX_TAG.SPAN}>{n}</Box><Box as={BOX_TAG.DIV}><Box as={BOX_TAG.H3}>{t}</Box><Box as={BOX_TAG.P}>{d}</Box></Box></Box>)}</Box>
        <Box as={BOX_TAG.DIV} className="rdr-admissions__highlight"><Sparkles size={20}/><Box as={BOX_TAG.SPAN}>{__("Every great school journey starts with a conversation. Our team is here to guide you through the next step.")}</Box></Box>
      </Box>
      <EnquiryForm/>
    </Box>
  </>
);

export { AcademicsPage };
export { CampusPage };
export { GalleryPage };
export { AdmissionsPage };
