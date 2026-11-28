import { _map } from '@Base/utilities/map';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Box , BOX_TAG} from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { AppLink } from '@Components/components/AppLink';
import { ArrowRight, ArrowUpRight, Check, ExternalLink, GraduationCap, HeartHandshake, Landmark, Sparkles } from '@Components/icons';
import { photographs } from '../../constants/school.content';
import { EnquiryLink } from '../siteActions';
import { Heading, ImageCard } from '../siteMedia';

const IntroPage = ({eyebrow,title,description,image}:{eyebrow:string;title:string;description:string;image?:string}) => (
  <Box as={BOX_TAG.SECTION} className="rdr-o1-page-hero">
    <Box as={BOX_TAG.DIV} className="rdr-container rdr-o1-page-hero__inner">
      <Box as={BOX_TAG.DIV} className="rdr-o1-page-hero__copy">
        <Box as={BOX_TAG.SPAN} className="rdr-eyebrow">{eyebrow}</Box>
        <Box as={BOX_TAG.H1}>{title}</Box>
        <Box as={BOX_TAG.P}>{description}</Box>
        <Box as={BOX_TAG.DIV} className="rdr-o1-page-hero__actions"><EnquiryLink className="rdr-btn rdr-btn--solid">{__("Enquire now ")}<ArrowUpRight size={16}/></EnquiryLink><AppLink to={SCHOOL_ROUTE.HOME} className="rdr-o1-page-hero__home">{__("Back to homepage ")}<ArrowRight size={16}/></AppLink></Box>
      </Box>
      {image&&<Box as={BOX_TAG.DIV} className="rdr-o1-page-hero__photo"><ImageCard src={image} alt={__("RDR School students and school life")} loading="eager"/><Box as={BOX_TAG.SPAN}>{__("RDR SCHOOL · HALUHERA")}</Box></Box>}
    </Box>
  </Box>
);

const AboutPage = () => (
  <>
    <IntroPage eyebrow="OUR STORY" title={__("A school with purpose, rooted in community.")} description={__("Rao Deena Ram Vidya Vihar was established in Haluhera in 2006, dedicated to thoughtful education and the all-round development of students.")} image={photographs[4].src}/>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container rdr-prose-grid">
      <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.SPAN} className="rdr-eyebrow">{__("OUR BEGINNINGS")}</Box><Box as={BOX_TAG.H2}>{__("Inspired by a legacy.")}</Box></Box>
      <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.P}>{__("The school is named in memory of Late Shri Rao Deena Ram Ji. According to the school's history, it was founded in April 2006 by Shri Lal Singh Yadav (Retired HES-II), to extend educational opportunities to the surrounding rural community.")}</Box><Box as={BOX_TAG.P}>{__("Its stated philosophy emphasizes respect, responsibility, fairness, caring and citizenship. These values continue to shape how the school's story is presented.")}</Box><ExternalAnchor href="https://rdrschool.com/our-school.html" target="_blank" rel="noreferrer" className="rdr-inline-link">{__("Read the original school history ")}<ExternalLink size={16}/></ExternalAnchor></Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-section--wash"><Box as={BOX_TAG.DIV} className="rdr-container rdr-values">
      {_map([{number:'01',title:'Learning',description:'Encouraging curiosity, creativity and independent thinking.'},{number:'02',title:'Character',description:'Building a foundation of responsibility, kindness and integrity.'},{number:'03',title:'Community',description:'Strengthening partnerships between families, teachers and students.'}], x=><Box as={BOX_TAG.DIV} key={x.number}><Box as={BOX_TAG.SPAN}>{x.number}</Box><Box as={BOX_TAG.H3}>{x.title}</Box><Box as={BOX_TAG.P}>{x.description}</Box></Box>)}
    </Box></Box>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container rdr-about-leaders">
      <Heading eyebrow="WORDS THAT GUIDE US" title={__("Meet the people and values behind RDR.")}/>
      <Box as={BOX_TAG.DIV} className="rdr-about-leaders__links">
        <AppLink to={SCHOOL_ROUTE.HISTORY}><Landmark size={22}/><Box as={BOX_TAG.STRONG}>{__("Our history")}</Box><Box as={BOX_TAG.SMALL}>{__("From Haluhera, since 2006")}</Box><ArrowUpRight size={16}/></AppLink>
        <AppLink to={SCHOOL_ROUTE.DIRECTOR}><HeartHandshake size={22}/><Box as={BOX_TAG.STRONG}>{__("Director's message")}</Box><Box as={BOX_TAG.SMALL}>{__("All-round development")}</Box><ArrowUpRight size={16}/></AppLink>
        <AppLink to={SCHOOL_ROUTE.PRINCIPAL}><GraduationCap size={22}/><Box as={BOX_TAG.STRONG}>{__("Principal's message")}</Box><Box as={BOX_TAG.SMALL}>{__("Learning with character")}</Box><ArrowUpRight size={16}/></AppLink>
        <AppLink to={SCHOOL_ROUTE.PHILOSOPHY}><Sparkles size={22}/><Box as={BOX_TAG.STRONG}>{__("Our philosophy")}</Box><Box as={BOX_TAG.SMALL}>{__("Broadening minds")}</Box><ArrowUpRight size={16}/></AppLink>
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container">
      <Heading eyebrow="LEADERSHIP & PHILOSOPHY" title={__("Guided by people who care.")} description={__("The school's historical leadership messages speak about holistic learning, kindness and preparing children for life beyond the classroom.")}/>
      <Box as={BOX_TAG.DIV} className="rdr-values rdr-values--leadership">
        <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.SPAN}>{__("FOUNDER & DIRECTOR · HISTORICAL PROFILE")}</Box><Box as={BOX_TAG.H3}>{__("Shri Lal Singh Yadav")}</Box><Box as={BOX_TAG.P}>{__("The original site credits Shri Lal Singh Yadav (Retired HES-II) with founding the school in April 2006. His published message emphasizes students' all-round development.")}</Box><ExternalAnchor className="rdr-inline-link" href="https://rdrschool.com/director-message.html" target="_blank" rel="noreferrer">{__("Original director's message ")}<ExternalLink size={16}/></ExternalAnchor></Box>
        <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.SPAN}>{__("PRINCIPAL · 2023 DISCLOSURE")}</Box><Box as={BOX_TAG.H3}>{__("Arvind Yadav")}</Box><Box as={BOX_TAG.P}>{__("Named as principal in the school's 2023 CBSE Appendix IX. Current office holder and leadership portrait must be verified before publication.")}</Box><ExternalAnchor className="rdr-inline-link" href="https://rdrschool.com/Principal-message.html" target="_blank" rel="noreferrer">{__("Original principal's message ")}<ExternalLink size={16}/></ExternalAnchor></Box>
        <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.SPAN}>{__("EDUCATIONAL PHILOSOPHY")}</Box><Box as={BOX_TAG.H3}>{__("Be a lamp unto yourself")}</Box><Box as={BOX_TAG.P}>{__("The legacy school describes its emblem as a lamp symbolizing personal understanding and responsibility. We retain that original identity alongside the RDR digital wordmark.")}</Box><ExternalAnchor className="rdr-inline-link" href="https://rdrschool.com/our-philosophy.html" target="_blank" rel="noreferrer">{__("School philosophy ")}<ExternalLink size={16}/></ExternalAnchor></Box>
      </Box>
    </Box>
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
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container rdr-story-page">
      <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.SPAN} className="rdr-eyebrow">{__("A CLOSER LOOK")}</Box><Box as={BOX_TAG.H2}>{__("Inside the RDR community.")}</Box><Box as={BOX_TAG.P}>{story.detail}</Box>
        <Box as={BOX_TAG.DIV} className="rdr-story-page__actions"><EnquiryLink className="rdr-btn rdr-btn--solid">{__("Ask the school ")}<ArrowUpRight size={16}/></EnquiryLink><ExternalAnchor href={story.src} target="_blank" rel="noreferrer" className="rdr-inline-link">{__("Original school message ")}<ExternalLink size={16}/></ExternalAnchor></Box>
      </Box>
      <Box as={BOX_TAG.DIV} className="rdr-story-page__notes">{_map(story.highlights, (line,i)=><Box as={BOX_TAG.DIV} key={line}><Box as={BOX_TAG.SPAN}>0{i+1}</Box><Box as={BOX_TAG.STRONG}>{line}</Box><Check size={16}/></Box>)}</Box>
    </Box>
  </>;
};

export { IntroPage };
export { AboutPage };
export { StoryPage };
