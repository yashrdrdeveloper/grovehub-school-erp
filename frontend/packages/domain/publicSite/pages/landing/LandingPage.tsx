import { SiteHeader, SiteFooter } from './components/siteChrome';
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
      <div className="rdr-site" data-testid="public-site-page" onClickCapture={handleInternalLink}>
        <SiteHeader/>
        <main id="main-content" data-testid="public-site-content"><CurrentPage/></main>
        <SiteFooter/>
        <PublicSheet panel={panel} onClose={closePanel} onNavigate={selectPanel} enquiry={<EnquiryForm compact/>}/>
      </div>
    </EnquiryActionContext.Provider>
  );
};
