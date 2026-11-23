import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../../../constants/school.routes';
import { AppLink } from '@Components/components/AppLink';
import { ArrowUpRight } from '@Components/icons';
import { facilities, photographs } from '../../../../constants/school.content';
import { Heading, ImageCard, facilityIcons } from '../../../siteMedia';
import { IntroPage } from '../../../aboutContent';

export const CampusPage = () => (
  <>
    <IntroPage
      eyebrow="CAMPUS & FACILITIES"
      title={__('Spaces for discovery and belonging.')}
      description={__(
        "From learning laboratories to playful outdoor moments, explore the facilities described on the school's original website.",
      )}
      image={photographs[2].src}
    />
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container">
      <Heading eyebrow="FACILITIES" title={__('Designed for a fuller school day.')} />
      <Box as={BOX_TAG.DIV} className="rdr-feature-grid">
        {_map(facilities, (x, i) => {
          const Icon = facilityIcons[i];
          return (
            <Box as={BOX_TAG.DIV} className="rdr-feature" key={x.title}>
              {Icon && <Icon size={29} strokeWidth={1.5} />}
              <Box as={BOX_TAG.H3}>{x.title}</Box>
              <Box as={BOX_TAG.P}>{x.description}</Box>
            </Box>
          );
        })}
      </Box>
      <Box as={BOX_TAG.DIV} className="rdr-campus-spotlight">
        <ImageCard src={photographs[2].src} alt={photographs[2].alt} />
        <Box as={BOX_TAG.DIV}>
          <Box as={BOX_TAG.SPAN} className="rdr-eyebrow">
            {__('CAMPUS LIFE')}
          </Box>
          <Box as={BOX_TAG.H3}>{__('Spaces to explore, memories to make.')}</Box>
          <Box as={BOX_TAG.P}>
            {__(
              'From early experiences to senior school, each day brings moments of curiosity, teamwork and discovery.',
            )}
          </Box>
          <AppLink className="rdr-inline-link" to={SCHOOL_ROUTE.GALLERY}>
            {__('See school life ')}
            <ArrowUpRight size={17} />
          </AppLink>
        </Box>
      </Box>
    </Box>
  </>
);
