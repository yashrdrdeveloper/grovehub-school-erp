import { Box, BOX_TAG } from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { AppLink } from '@Components/components/AppLink';
import { ArrowRight, ArrowUpRight } from '@Components/icons';
import { photographs } from '../../constants/school.content';
import { EnquiryLink } from '../siteActions';
import { ImageCard } from '../siteMedia';

const HomePage = () => (
  <Box as={BOX_TAG.DIV} className="rdr-option-one">
    <Box as={BOX_TAG.SECTION} className="rdr-o1-hero">
      <Box as={BOX_TAG.DIV} className="rdr-container rdr-o1-hero__grid">
        <Box as={BOX_TAG.DIV} className="rdr-o1-hero__copy">
          <Box as={BOX_TAG.SPAN} className="rdr-eyebrow">
            A COMMUNITY OF LEARNERS · SINCE 2006
          </Box>
          <Box as={BOX_TAG.H1}>
            Rooted in Haluhera.<em>Ready for every tomorrow.</em>
          </Box>
          <Box as={BOX_TAG.P}>
            At Rao Deena Ram Vidya Vihar, learning goes beyond the classroom. We bring together
            academics, character and a caring school community so every child has room to grow.
          </Box>
          <Box as={BOX_TAG.DIV} className="rdr-o1-hero__actions">
            <EnquiryLink className="rdr-btn rdr-btn--solid">
              Enquire for admission <ArrowUpRight size={17} />
            </EnquiryLink>
            <AppLink className="rdr-btn rdr-o1-secondary" to={SCHOOL_ROUTE.ABOUT}>
              Discover our school <ArrowRight size={17} />
            </AppLink>
          </Box>
          <Box as={BOX_TAG.DIV} className="rdr-o1-hero__proof">
            <Box as={BOX_TAG.SPAN}>CBSE affiliated</Box>
            <Box as={BOX_TAG.SPAN}>Senior Secondary</Box>
            <Box as={BOX_TAG.SPAN}>Rewari, Haryana</Box>
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className="rdr-o1-hero__images">
          <ImageCard
            className="rdr-o1-hero__main-photo"
            src={photographs[4].src}
            alt={photographs[4].alt}
            position={photographs[4].position}
            loading="eager"
          >
            <Box as={BOX_TAG.SPAN} className="rdr-o1-hero__photo-caption">
              A school full of familiar faces.
            </Box>
          </ImageCard>
          <Box as={BOX_TAG.DIV} className="rdr-o1-hero__image-side">
            <ImageCard
              className="rdr-o1-hero__side-photo"
              src={photographs[0].src}
              alt={photographs[0].alt}
              position={photographs[0].position}
              loading="eager"
            />
            <Box as={BOX_TAG.DIV} className="rdr-o1-hero__note">
              <Box as={BOX_TAG.STRONG}>Learning with purpose.</Box>
              <Box as={BOX_TAG.SMALL}>Growing together since 2006.</Box>
            </Box>
          </Box>
        </Box>
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className="rdr-o1-life rdr-container">
      <Box as={BOX_TAG.DIV} className="rdr-o1-life__heading">
        <Box as={BOX_TAG.DIV}>
          <Box as={BOX_TAG.SPAN} className="rdr-eyebrow">
            LIFE AT RDR
          </Box>
          <Box as={BOX_TAG.H2}>
            Every child has a story.
            <br />
            Every day has a purpose.
          </Box>
        </Box>
        <Box as={BOX_TAG.P}>
          From classrooms to celebrations, these are the experiences that shape school life.
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className="rdr-o1-life__grid">
        <AppLink to={SCHOOL_ROUTE.GALLERY} className="rdr-o1-life__card">
          <ImageCard
            src={photographs[0].src}
            alt={photographs[0].alt}
            position={photographs[0].position}
          />
          <Box as={BOX_TAG.SPAN}>ACHIEVEMENT</Box>
          <Box as={BOX_TAG.H3}>Celebrating every milestone.</Box>
          <Box as={BOX_TAG.P}>Recognising effort, confidence and achievement.</Box>
        </AppLink>
        <AppLink to={SCHOOL_ROUTE.GALLERY} className="rdr-o1-life__card">
          <ImageCard
            src={photographs[1].src}
            alt={photographs[1].alt}
            position={photographs[1].position}
          />
          <Box as={BOX_TAG.SPAN}>EXPRESSION</Box>
          <Box as={BOX_TAG.H3}>More than textbooks.</Box>
          <Box as={BOX_TAG.P}>Music, movement and the joy of participating.</Box>
        </AppLink>
        <AppLink to={SCHOOL_ROUTE.GALLERY} className="rdr-o1-life__card">
          <ImageCard
            src={photographs[2].src}
            alt={photographs[2].alt}
            position={photographs[2].position}
          />
          <Box as={BOX_TAG.SPAN}>BELONGING</Box>
          <Box as={BOX_TAG.H3}>Growing with confidence.</Box>
          <Box as={BOX_TAG.P}>Creating space for curiosity and friendships.</Box>
        </AppLink>
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className="rdr-o1-cta">
      <Box as={BOX_TAG.DIV} className="rdr-container rdr-o1-cta__inner">
        <Box as={BOX_TAG.DIV}>
          <Box as={BOX_TAG.H2}>Come see what makes RDR special.</Box>
          <Box as={BOX_TAG.P}>Questions about admissions, classes or a campus visit?</Box>
        </Box>
        <EnquiryLink className="rdr-btn rdr-btn--light">
          Start an enquiry <ArrowUpRight size={17} />
        </EnquiryLink>
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className="rdr-o1-disclosure rdr-container">
      <Box as={BOX_TAG.SPAN}>Looking for school documents?</Box>
      <AppLink to={SCHOOL_ROUTE.DISCLOSURE}>
        Mandatory Public Disclosure <ArrowUpRight size={16} />
      </AppLink>
    </Box>
  </Box>
);

export { HomePage };
