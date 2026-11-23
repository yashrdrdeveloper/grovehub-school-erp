import { cn } from '@Base/utilities/cn';
import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './HomeContent.module.css';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { AppLink } from '@Components/components/AppLink';
import { ArrowRight, ArrowUpRight } from '@Components/icons';
import { photographs } from '../../constants/school.content';
import { EnquiryLink } from '../siteActions';
import { ImageCard } from '../siteMedia';

const HomePage = () => (
  <Box as={BOX_TAG.DIV} className={styles['rdr-option-one']}>
    <Box as={BOX_TAG.SECTION} className={styles['rdr-o1-hero']}>
      <Box as={BOX_TAG.DIV} className={cn('rdr-container', styles['rdr-o1-hero__grid'])}>
        <Box as={BOX_TAG.DIV} className={styles['rdr-o1-hero__copy']}>
          <Box as={BOX_TAG.SPAN} className="rdr-eyebrow">
            {__('A COMMUNITY OF LEARNERS · SINCE 2006')}
          </Box>
          <Box as={BOX_TAG.H1}>
            {__('Rooted in Haluhera.')}
            <em>{__('Ready for every tomorrow.')}</em>
          </Box>
          <Box as={BOX_TAG.P}>
            {__(
              'At Rao Deena Ram Vidya Vihar, learning goes beyond the classroom. We bring together academics, character and a caring school community so every child has room to grow.',
            )}
          </Box>
          <Box as={BOX_TAG.DIV} className={styles['rdr-o1-hero__actions']}>
            <EnquiryLink className="rdr-btn rdr-btn--solid">
              {__('Enquire for admission ')}
              <ArrowUpRight size={17} />
            </EnquiryLink>
            <AppLink className={cn('rdr-btn', styles['rdr-o1-secondary'])} to={SCHOOL_ROUTE.ABOUT}>
              {__('Discover our school ')}
              <ArrowRight size={17} />
            </AppLink>
          </Box>
          <Box as={BOX_TAG.DIV} className={styles['rdr-o1-hero__proof']}>
            <Box as={BOX_TAG.SPAN}>{__('CBSE affiliated')}</Box>
            <Box as={BOX_TAG.SPAN}>{__('Senior Secondary')}</Box>
            <Box as={BOX_TAG.SPAN}>{__('Rewari, Haryana')}</Box>
          </Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles['rdr-o1-hero__images']}>
          <ImageCard
            className={styles['rdr-o1-hero__main-photo']}
            src={photographs[4].src}
            alt={photographs[4].alt}
            position={photographs[4].position}
            loading="eager"
          >
            <Box as={BOX_TAG.SPAN} className={styles['rdr-o1-hero__photo-caption']}>
              {__('A school full of familiar faces.')}
            </Box>
          </ImageCard>
          <Box as={BOX_TAG.DIV} className={styles['rdr-o1-hero__image-side']}>
            <ImageCard
              className={styles['rdr-o1-hero__side-photo']}
              src={photographs[0].src}
              alt={photographs[0].alt}
              position={photographs[0].position}
              loading="eager"
            />
            <Box as={BOX_TAG.DIV} className={styles['rdr-o1-hero__note']}>
              <Box as={BOX_TAG.STRONG}>{__('Learning with purpose.')}</Box>
              <Box as={BOX_TAG.SMALL}>{__('Growing together since 2006.')}</Box>
            </Box>
          </Box>
        </Box>
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className={cn(styles['rdr-o1-life'], 'rdr-container')}>
      <Box as={BOX_TAG.DIV} className={styles['rdr-o1-life__heading']}>
        <Box as={BOX_TAG.DIV}>
          <Box as={BOX_TAG.SPAN} className="rdr-eyebrow">
            {__('LIFE AT RDR')}
          </Box>
          <Box as={BOX_TAG.H2}>
            {__('Every child has a story.')}
            <br />
            {__('Every day has a purpose.')}
          </Box>
        </Box>
        <Box as={BOX_TAG.P}>
          {__('From classrooms to celebrations, these are the experiences that shape school life.')}
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className={styles['rdr-o1-life__grid']}>
        <AppLink to={SCHOOL_ROUTE.GALLERY} className={styles['rdr-o1-life__card']}>
          <ImageCard
            src={photographs[0].src}
            alt={photographs[0].alt}
            position={photographs[0].position}
          />
          <Box as={BOX_TAG.SPAN}>{__('ACHIEVEMENT')}</Box>
          <Box as={BOX_TAG.H3}>{__('Celebrating every milestone.')}</Box>
          <Box as={BOX_TAG.P}>{__('Recognising effort, confidence and achievement.')}</Box>
        </AppLink>
        <AppLink to={SCHOOL_ROUTE.GALLERY} className={styles['rdr-o1-life__card']}>
          <ImageCard
            src={photographs[1].src}
            alt={photographs[1].alt}
            position={photographs[1].position}
          />
          <Box as={BOX_TAG.SPAN}>{__('EXPRESSION')}</Box>
          <Box as={BOX_TAG.H3}>{__('More than textbooks.')}</Box>
          <Box as={BOX_TAG.P}>{__('Music, movement and the joy of participating.')}</Box>
        </AppLink>
        <AppLink to={SCHOOL_ROUTE.GALLERY} className={styles['rdr-o1-life__card']}>
          <ImageCard
            src={photographs[2].src}
            alt={photographs[2].alt}
            position={photographs[2].position}
          />
          <Box as={BOX_TAG.SPAN}>{__('BELONGING')}</Box>
          <Box as={BOX_TAG.H3}>{__('Growing with confidence.')}</Box>
          <Box as={BOX_TAG.P}>{__('Creating space for curiosity and friendships.')}</Box>
        </AppLink>
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className={styles['rdr-o1-cta']}>
      <Box as={BOX_TAG.DIV} className={cn('rdr-container', styles['rdr-o1-cta__inner'])}>
        <Box as={BOX_TAG.DIV}>
          <Box as={BOX_TAG.H2}>{__('Come see what makes RDR special.')}</Box>
          <Box as={BOX_TAG.P}>{__('Questions about admissions, classes or a campus visit?')}</Box>
        </Box>
        <EnquiryLink className="rdr-btn rdr-btn--light">
          {__('Start an enquiry ')}
          <ArrowUpRight size={17} />
        </EnquiryLink>
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className={cn(styles['rdr-o1-disclosure'], 'rdr-container')}>
      <Box as={BOX_TAG.SPAN}>{__('Looking for school documents?')}</Box>
      <AppLink to={SCHOOL_ROUTE.DISCLOSURE}>
        {__('Mandatory Public Disclosure ')}
        <ArrowUpRight size={16} />
      </AppLink>
    </Box>
  </Box>
);

export { HomePage };
