import { _map } from '@Base/utilities/map';
import { Box } from '@Components/components/Box';
import { ArrowLeft, LockKeyhole } from '@Components/icons';
import { Link } from '@Components/router';
import { APP_ROUTE } from '@Business/app/routes';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme } from '@Components/hooks/useWorkspaceTheme';
import { LOGIN_LAYOUT_APPEARANCES } from './constants/loginLayout.appearances';
import { LOGIN_LAYOUT_TEST_IDS } from './constants/loginLayout.testIds';
import { getLoginLayoutClassName } from './helpers/loginLayout.getClassName';
import { renderLoginLayoutAppearanceOption } from './renderer/loginLayout.appearanceOption';
import type { LoginLayoutComponent } from './types/loginLayout.props';

import { AuthAmbientGraphic } from '../authAmbientGraphic';
import { useAuthViewportFit } from '../../hooks/useAuthViewportFit';


// One stylesheet owns the entire auth layout. Legacy variants are intentionally not loaded.
import './styles/optionOneAuth.css';

/** Shared, card-free shell for sign-in and every recovery step. */
export const LoginLayout: LoginLayoutComponent = ({ children, title, description, mobileDescription, minimalMobile = false, variant = 'login', eyebrow = 'YOUR SCHOOL WORKSPACE' }) => {
  const { preference, setPreference } = useWorkspaceTheme();
  const { contentRef, toolbarRef, mainRef } = useAuthViewportFit();


  return (
    <Box as="main" className={getLoginLayoutClassName(variant, minimalMobile)} data-testid={LOGIN_LAYOUT_TEST_IDS.ROOT}>
      <Box as="div" className="rdr-login-v2__frame">
        <Box as="section" className="rdr-login-v2__visual" aria-label="RDR School branded gradient background">
          <Box as="div" className="rdr-login-v2__mobile-header">
            <Box as="div" className="rdr-login-v2__mobile-identity" aria-label="RDR School">
              <SchoolMonogram/>
              <Box as="span" className="rdr-login-v2__mobile-school-name">RDR School<Box as="small">Your school workspace</Box></Box>
            </Box>
          </Box>
          <Box as="div" className="rdr-login-v2__visual-brand">
            <SchoolMonogram/>
            <Box as="span"><Box as="strong">RDR School</Box><Box as="small">Rao Deena Ram Vidya Vihar</Box></Box>
          </Box>
          <AuthAmbientGraphic/>
          <Box as="div" className="rdr-login-v2__visual-copy">
            <Box as="span" className="rdr-login-v2__visual-kicker"><Box as="span"/> THE RDR SCHOOL WORKSPACE</Box>
            <Box as="h2">For every school day.<br/><em>And everyone behind it.</em></Box>
            <Box as="p">A connected space for the people who make learning happen.</Box>
            <Box as="div" className="rdr-login-v2__visual-footer"><Box as="span">LEARN</Box><i aria-hidden="true"/><Box as="span">ORGANISE</Box><i aria-hidden="true"/><Box as="span">GROW</Box></Box>
          </Box>
        </Box>
        <Box as="section" ref={contentRef} className="rdr-login-v2__content" data-testid={LOGIN_LAYOUT_TEST_IDS.CONTENT} aria-labelledby="rdr-auth-title">
          <Box as="div" ref={toolbarRef} className="rdr-login-v2__toolbar" data-testid={LOGIN_LAYOUT_TEST_IDS.TOOLBAR}>
            <Link to={APP_ROUTE.HOME} className="rdr-login-v2__back" aria-label="Back to RDR School homepage">
              <ArrowLeft size={19}/><Box as="span">Back to school</Box>
            </Link>
            <Box as="div" className="rdr-login-v2__appearance" role="group" aria-label="Appearance">
              {_map(LOGIN_LAYOUT_APPEARANCES, option => renderLoginLayoutAppearanceOption({ ...option, preference, setPreference }))}
            </Box>
          </Box>
          <Box as="div" ref={mainRef} className="rdr-login-v2__main" data-testid={LOGIN_LAYOUT_TEST_IDS.MAIN}>
            <Box as="div" className="rdr-login-v2__heading">
              <Box as="span" className="rdr-login-v2__eyebrow">{eyebrow}</Box>
              <Box as="h1" id="rdr-auth-title" data-testid={LOGIN_LAYOUT_TEST_IDS.TITLE}>{title}</Box>
              <Box as="p" className="rdr-login-v2__description-desktop">{description}</Box>
              <Box as="p" className="rdr-login-v2__description-mobile">{mobileDescription ?? description}</Box>
            </Box>
            {children}
            <Box as="div" className="rdr-login-v2__closing" aria-label="Secure RDR School workspace">
              <LockKeyhole size={13} strokeWidth={1.8} aria-hidden="true"/>
              <Box as="span">A secure workspace for RDR School</Box>
            </Box>
            <Box as="div" className="rdr-o1-auth-footer"><Box as="span">© RDR School</Box><Box as="span">Haluhera · Rewari</Box></Box>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};
