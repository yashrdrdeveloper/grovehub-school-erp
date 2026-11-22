import { Box } from '@Components/components/Box';
import { SiteHeader, SiteFooter } from './components/siteChrome';
import { LANDING_TEST_IDS } from './constants/landing.testIds';
import { EnquiryActionContext } from './components/siteActions';
import { EnquiryForm } from './components/enquiryForm';
import { PublicSheet } from './components/publicSheet';
import type { LandingPageModel } from './hooks/useLandingPage/types/useLandingPage';

import './styles/landingPage.css';
import './styles/optionOne.css';
import './styles/optionOnePages.css';

/** Public route presentation only; there is no public workspace/layout wrapper. */
export const LandingPage = ({ model }: { readonly model: LandingPageModel }) => {
  const { CurrentPage, panel, closePanel, selectPanel, handleInternalLink, openEnquiry } = model;
  return (
    <EnquiryActionContext.Provider value={openEnquiry}>
      <Box as="div" className="rdr-site" data-testid={LANDING_TEST_IDS.PAGE} onClickCapture={handleInternalLink}>
        <SiteHeader/>
        <Box as="main" id="main-content" data-testid={LANDING_TEST_IDS.CONTENT}><CurrentPage/></Box>
        <SiteFooter/>
        <PublicSheet panel={panel} onClose={closePanel} onNavigate={selectPanel} enquiry={<EnquiryForm compact/>}/>
      </Box>
    </EnquiryActionContext.Provider>
  );
};
