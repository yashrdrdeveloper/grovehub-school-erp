import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';
import { ArrowLeft, LockKeyhole } from '@Components/icons';
import { AppLink } from '@Components/components/AppLink';
import { APP_ROUTE } from '@Business/app/routes';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme } from '@Components/hooks/useWorkspaceTheme';
import { LOGIN_LAYOUT_APPEARANCES } from './constants/loginLayout.appearances';
import { LOGIN_LAYOUT_TEST_IDS } from './constants/loginLayout.testIds';
import { LoginLayoutVariant } from './constants/loginLayout.variant';
import { LOGIN_LAYOUT_COPY } from './constants/loginLayout.copy';
import { getLoginLayoutClassName } from './helpers/loginLayout.getClassName';
import { renderLoginLayoutAppearanceOption } from './renderer/loginLayout.appearanceOption';
import type { LoginLayoutComponent } from './types/loginLayout.props';

import { AuthAmbientGraphic } from '../authAmbientGraphic';
import { useAuthViewportFit } from '../../hooks/useAuthViewportFit';

// One stylesheet owns the entire auth layout. Legacy variants are intentionally not loaded.
import './styles/optionOneAuth.css';

/** Shared, card-free shell for sign-in and every recovery step. */
export const LoginLayout: LoginLayoutComponent = ({
  children,
  title,
  description,
  mobileDescription,
  minimalMobile = false,
  variant = LoginLayoutVariant.LOGIN,
  eyebrow = LOGIN_LAYOUT_COPY.DEFAULT_EYEBROW,
}) => {
  const { preference, setPreference } = useWorkspaceTheme();
  const { contentRef, toolbarRef, mainRef } = useAuthViewportFit();

  return (
    <Box
      as={BOX_TAG.MAIN}
      className={getLoginLayoutClassName(variant, minimalMobile)}
      data-testid={LOGIN_LAYOUT_TEST_IDS.ROOT}
    >
      <Box as={BOX_TAG.DIV} className="rdr-login-v2__frame">
        <Box
          as={BOX_TAG.SECTION}
          className="rdr-login-v2__visual"
          aria-label="RDR School branded gradient background"
        >
          <Box as={BOX_TAG.DIV} className="rdr-login-v2__mobile-header">
            <Box as={BOX_TAG.DIV} className="rdr-login-v2__mobile-identity" aria-label="RDR School">
              <SchoolMonogram />
              <Box as={BOX_TAG.SPAN} className="rdr-login-v2__mobile-school-name">
                RDR School<Box as={BOX_TAG.SMALL}>Your school workspace</Box>
              </Box>
            </Box>
          </Box>
          <Box as={BOX_TAG.DIV} className="rdr-login-v2__visual-brand">
            <SchoolMonogram />
            <Box as={BOX_TAG.SPAN}>
              <Box as={BOX_TAG.STRONG}>RDR School</Box>
              <Box as={BOX_TAG.SMALL}>Rao Deena Ram Vidya Vihar</Box>
            </Box>
          </Box>
          <AuthAmbientGraphic />
          <Box as={BOX_TAG.DIV} className="rdr-login-v2__visual-copy">
            <Box as={BOX_TAG.SPAN} className="rdr-login-v2__visual-kicker">
              <Box as={BOX_TAG.SPAN} /> THE RDR SCHOOL WORKSPACE
            </Box>
            <Box as={BOX_TAG.H2}>
              For every school day.
              <br />
              <em>And everyone behind it.</em>
            </Box>
            <Box as={BOX_TAG.P}>A connected space for the people who make learning happen.</Box>
            <Box as={BOX_TAG.DIV} className="rdr-login-v2__visual-footer">
              <Box as={BOX_TAG.SPAN}>LEARN</Box>
              <i aria-hidden="true" />
              <Box as={BOX_TAG.SPAN}>ORGANISE</Box>
              <i aria-hidden="true" />
              <Box as={BOX_TAG.SPAN}>GROW</Box>
            </Box>
          </Box>
        </Box>
        <Box
          as={BOX_TAG.SECTION}
          ref={contentRef}
          className="rdr-login-v2__content"
          data-testid={LOGIN_LAYOUT_TEST_IDS.CONTENT}
          aria-labelledby="rdr-auth-title"
        >
          <Box
            as={BOX_TAG.DIV}
            ref={toolbarRef}
            className="rdr-login-v2__toolbar"
            data-testid={LOGIN_LAYOUT_TEST_IDS.TOOLBAR}
          >
            <AppLink
              to={APP_ROUTE.HOME}
              className="rdr-login-v2__back"
              aria-label="Back to RDR School homepage"
            >
              <ArrowLeft size={19} />
              <Box as={BOX_TAG.SPAN}>Back to school</Box>
            </AppLink>
            <Box
              as={BOX_TAG.DIV}
              className="rdr-login-v2__appearance"
              role="group"
              aria-label="Appearance"
            >
              {_map(LOGIN_LAYOUT_APPEARANCES, (option) =>
                renderLoginLayoutAppearanceOption({ ...option, preference, setPreference }),
              )}
            </Box>
          </Box>
          <Box
            as={BOX_TAG.DIV}
            ref={mainRef}
            className="rdr-login-v2__main"
            data-testid={LOGIN_LAYOUT_TEST_IDS.MAIN}
          >
            <Box as={BOX_TAG.DIV} className="rdr-login-v2__heading">
              <Box as={BOX_TAG.SPAN} className="rdr-login-v2__eyebrow">
                {eyebrow}
              </Box>
              <Box as={BOX_TAG.H1} id="rdr-auth-title" data-testid={LOGIN_LAYOUT_TEST_IDS.TITLE}>
                {title}
              </Box>
              <Box as={BOX_TAG.P} className="rdr-login-v2__description-desktop">
                {description}
              </Box>
              <Box as={BOX_TAG.P} className="rdr-login-v2__description-mobile">
                {mobileDescription ?? description}
              </Box>
            </Box>
            {children}
            <Box
              as={BOX_TAG.DIV}
              className="rdr-login-v2__closing"
              aria-label="Secure RDR School workspace"
            >
              <LockKeyhole size={13} strokeWidth={1.8} aria-hidden="true" />
              <Box as={BOX_TAG.SPAN}>A secure workspace for RDR School</Box>
            </Box>
            <Box as={BOX_TAG.DIV} className="rdr-o1-auth-footer">
              <Box as={BOX_TAG.SPAN}>© RDR School</Box>
              <Box as={BOX_TAG.SPAN}>Haluhera · Rewari</Box>
            </Box>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};
