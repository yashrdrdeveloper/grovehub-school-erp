import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Box } from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { Link } from '@Components/router';
import { ArrowRight, ArrowUpRight, Check, ExternalLink, GraduationCap, HeartHandshake, Landmark, Sparkles } from '@Components/icons';
import { photographs } from '../../constants/school.content';
import { EnquiryLink } from '../siteActions';
import { Heading, ImageCard } from '../siteMedia';

const IntroPage = ({eyebrow,title,description,image}:{eyebrow:string;title:string;description:string;image?:string}) => (
  <Box as="section" className="rdr-o1-page-hero">
    <Box as="div" className="rdr-container rdr-o1-page-hero__inner">
      <Box as="div" className="rdr-o1-page-hero__copy">
        <Box as="span" className="rdr-eyebrow">{eyebrow}</Box>
        <Box as="h1">{title}</Box>
        <Box as="p">{description}</Box>
        <Box as="div" className="rdr-o1-page-hero__actions"><EnquiryLink className="rdr-btn rdr-btn--solid">Enquire now <ArrowUpRight size={16}/></EnquiryLink><Link to={SCHOOL_ROUTE.HOME} className="rdr-o1-page-hero__home">Back to homepage <ArrowRight size={16}/></Link></Box>
      </Box>
      {image&&<Box as="div" className="rdr-o1-page-hero__photo"><ImageCard src={image} alt="RDR School students and school life" loading="eager"/><Box as="span">RDR SCHOOL · HALUHERA</Box></Box>}
    </Box>
  </Box>
);

const AboutPage = () => (
  <>
    <IntroPage eyebrow="OUR STORY" title="A school with purpose, rooted in community." description="Rao Deena Ram Vidya Vihar was established in Haluhera in 2006, dedicated to thoughtful education and the all-round development of students." image={photographs[4].src}/>
    <Box as="section" className="rdr-section rdr-container rdr-prose-grid">
      <Box as="div"><Box as="span" className="rdr-eyebrow">OUR BEGINNINGS</Box><Box as="h2">Inspired by a legacy.</Box></Box>
      <Box as="div"><Box as="p">The school is named in memory of Late Shri Rao Deena Ram Ji. According to the school's history, it was founded in April 2006 by Shri Lal Singh Yadav (Retired HES-II), to extend educational opportunities to the surrounding rural community.</Box><Box as="p">Its stated philosophy emphasizes respect, responsibility, fairness, caring and citizenship. These values continue to shape how the school's story is presented.</Box><ExternalAnchor href="https://rdrschool.com/our-school.html" target="_blank" rel="noreferrer" className="rdr-inline-link">Read the original school history <ExternalLink size={16}/></ExternalAnchor></Box>
    </Box>
    <Box as="section" className="rdr-section rdr-section--wash"><Box as="div" className="rdr-container rdr-values">
      {[{number:'01',title:'Learning',description:'Encouraging curiosity, creativity and independent thinking.'},{number:'02',title:'Character',description:'Building a foundation of responsibility, kindness and integrity.'},{number:'03',title:'Community',description:'Strengthening partnerships between families, teachers and students.'}].map(x=><Box as="div" key={x.number}><Box as="span">{x.number}</Box><Box as="h3">{x.title}</Box><Box as="p">{x.description}</Box></Box>)}
    </Box></Box>
    <Box as="section" className="rdr-section rdr-container rdr-about-leaders">
      <Heading eyebrow="WORDS THAT GUIDE US" title="Meet the people and values behind RDR."/>
      <Box as="div" className="rdr-about-leaders__links">
        <Link to={SCHOOL_ROUTE.HISTORY}><Landmark size={22}/><Box as="strong">Our history</Box><Box as="small">From Haluhera, since 2006</Box><ArrowUpRight size={16}/></Link>
        <Link to={SCHOOL_ROUTE.DIRECTOR}><HeartHandshake size={22}/><Box as="strong">Director's message</Box><Box as="small">All-round development</Box><ArrowUpRight size={16}/></Link>
        <Link to={SCHOOL_ROUTE.PRINCIPAL}><GraduationCap size={22}/><Box as="strong">Principal's message</Box><Box as="small">Learning with character</Box><ArrowUpRight size={16}/></Link>
        <Link to={SCHOOL_ROUTE.PHILOSOPHY}><Sparkles size={22}/><Box as="strong">Our philosophy</Box><Box as="small">Broadening minds</Box><ArrowUpRight size={16}/></Link>
      </Box>
    </Box>
    <Box as="section" className="rdr-section rdr-container">
      <Heading eyebrow="LEADERSHIP & PHILOSOPHY" title="Guided by people who care." description="The school's historical leadership messages speak about holistic learning, kindness and preparing children for life beyond the classroom."/>
      <Box as="div" className="rdr-values rdr-values--leadership">
        <Box as="div"><Box as="span">FOUNDER & DIRECTOR · HISTORICAL PROFILE</Box><Box as="h3">Shri Lal Singh Yadav</Box><Box as="p">The original site credits Shri Lal Singh Yadav (Retired HES-II) with founding the school in April 2006. His published message emphasizes students' all-round development.</Box><ExternalAnchor className="rdr-inline-link" href="https://rdrschool.com/director-message.html" target="_blank" rel="noreferrer">Original director's message <ExternalLink size={16}/></ExternalAnchor></Box>
        <Box as="div"><Box as="span">PRINCIPAL · 2023 DISCLOSURE</Box><Box as="h3">Arvind Yadav</Box><Box as="p">Named as principal in the school's 2023 CBSE Appendix IX. Current office holder and leadership portrait must be verified before publication.</Box><ExternalAnchor className="rdr-inline-link" href="https://rdrschool.com/Principal-message.html" target="_blank" rel="noreferrer">Original principal's message <ExternalLink size={16}/></ExternalAnchor></Box>
        <Box as="div"><Box as="span">EDUCATIONAL PHILOSOPHY</Box><Box as="h3">Be a lamp unto yourself</Box><Box as="p">The legacy school describes its emblem as a lamp symbolizing personal understanding and responsibility. We retain that original identity alongside the RDR digital wordmark.</Box><ExternalAnchor className="rdr-inline-link" href="https://rdrschool.com/our-philosophy.html" target="_blank" rel="noreferrer">School philosophy <ExternalLink size={16}/></ExternalAnchor></Box>
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
    <Box as="section" className="rdr-section rdr-container rdr-story-page">
      <Box as="div"><Box as="span" className="rdr-eyebrow">A CLOSER LOOK</Box><Box as="h2">Inside the RDR community.</Box><Box as="p">{story.detail}</Box>
        <Box as="div" className="rdr-story-page__actions"><EnquiryLink className="rdr-btn rdr-btn--solid">Ask the school <ArrowUpRight size={16}/></EnquiryLink><ExternalAnchor href={story.src} target="_blank" rel="noreferrer" className="rdr-inline-link">Original school message <ExternalLink size={16}/></ExternalAnchor></Box>
      </Box>
      <Box as="div" className="rdr-story-page__notes">{story.highlights.map((line,i)=><Box as="div" key={line}><Box as="span">0{i+1}</Box><Box as="strong">{line}</Box><Check size={16}/></Box>)}</Box>
    </Box>
  </>;
};

export { IntroPage };
export { AboutPage };
export { StoryPage };
