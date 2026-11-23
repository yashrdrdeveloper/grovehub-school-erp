import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';
import { Sparkles } from '@Components/icons';
import { photographs } from '../../../../constants/school.content';
import { Heading } from '../../../siteMedia';
import { IntroPage } from '../../../aboutContent';
import { EnquiryForm } from '../../../enquiryForm';

export const AdmissionsPage = () => (
  <>
    <IntroPage
      eyebrow="JOIN OUR COMMUNITY"
      title={__("Your child's journey starts here.")}
      description={__(
        'Exploring school options? We can help with admissions, school visits, fees, transport and academic questions.',
      )}
      image={photographs[3].src}
    />
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container rdr-admissions-layout">
      <Box as={BOX_TAG.DIV}>
        <Heading
          eyebrow="HOW TO BEGIN"
          title={__('An easier first step.')}
          description={__(
            "The school's historical admission procedure asks families to contact the school office and collect an application. Current eligibility, seat availability and the academic-year schedule should be verified directly.",
          )}
        />
        <Box as={BOX_TAG.DIV} className="rdr-steps">
          {_map(
            [
              ['01', 'Share an enquiry', 'Tell us what you would like to know.'],
              ['02', 'Talk to the school', 'Discuss class availability, visits and requirements.'],
              [
                '03',
                'Follow the admission process',
                'Get current forms and guidance directly from the school.',
              ],
            ],
            ([n, t, d]) => (
              <Box as={BOX_TAG.DIV} key={n}>
                <Box as={BOX_TAG.SPAN}>{n}</Box>
                <Box as={BOX_TAG.DIV}>
                  <Box as={BOX_TAG.H3}>{t}</Box>
                  <Box as={BOX_TAG.P}>{d}</Box>
                </Box>
              </Box>
            ),
          )}
        </Box>
        <Box as={BOX_TAG.DIV} className="rdr-admissions__highlight">
          <Sparkles size={20} />
          <Box as={BOX_TAG.SPAN}>
            {__(
              'Every great school journey starts with a conversation. Our team is here to guide you through the next step.',
            )}
          </Box>
        </Box>
      </Box>
      <EnquiryForm />
    </Box>
  </>
);
