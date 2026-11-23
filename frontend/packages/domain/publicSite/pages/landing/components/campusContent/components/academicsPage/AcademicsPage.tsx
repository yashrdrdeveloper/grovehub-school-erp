import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';
import { BookOpen } from '@Components/icons';
import { facilities, photographs } from '../../../../constants/school.content';
import { Heading } from '../../../siteMedia';
import { IntroPage } from '../../../aboutContent';

export const AcademicsPage = () => (
  <>
    <IntroPage
      eyebrow="ACADEMICS"
      title={__('Learning that shapes a lifetime.')}
      description={__(
        'A broad educational experience focused on understanding, confidence, creativity and practical learning.',
      )}
      image={photographs[0].src}
    />
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container">
      <Heading
        eyebrow="LEARNING JOURNEY"
        title={__('Room to think. Space to discover.')}
        description={__(
          "From the earliest years through senior secondary school, our focus is on building strong foundations and supporting each learner's growth.",
        )}
      />
      <Box as={BOX_TAG.DIV} className="rdr-feature-grid">
        {_map(
          [
            ['Foundational Years', 'Exploration, play and early learning experiences.'],
            ['Primary Years', 'Building literacy, numeracy and confidence.'],
            ['Middle School', 'Developing independence and deeper understanding.'],
            [
              'Secondary & Senior Secondary',
              'Preparing students for advanced study and future pathways.',
            ],
          ],
          ([title, description], i) => (
            <Box as={BOX_TAG.DIV} key={title} className="rdr-feature">
              <Box as={BOX_TAG.SPAN} className="rdr-feature__number">
                0{i + 1}
              </Box>
              <BookOpen size={28} />
              <Box as={BOX_TAG.H3}>{title}</Box>
              <Box as={BOX_TAG.P}>{description}</Box>
            </Box>
          ),
        )}
      </Box>
    </Box>
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-section--wash">
      <Box as={BOX_TAG.DIV} className="rdr-container rdr-prose-grid">
        <Box as={BOX_TAG.DIV}>
          <Heading eyebrow="BEYOND THE CLASSROOM" title={__('Learning in action.')} />
        </Box>
        <Box as={BOX_TAG.P}>
          {__(
            'Our historical school information describes computer education, science facilities, reading resources, sports, cultural programmes and student activities. For current class offerings, streams, subject combinations and curricula, please contact the school office.',
          )}
        </Box>
      </Box>
    </Box>
  </>
);
