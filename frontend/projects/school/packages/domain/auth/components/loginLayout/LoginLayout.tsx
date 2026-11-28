import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';
import { ArrowLeft } from '@Components/icons';
import { AppLink } from '@Components/components/AppLink';
import { APP_ROUTE } from '@Business/app/routes';

import { useWorkspaceTheme } from '@Components/hooks/useWorkspaceTheme';
import { SCHOOL_THEME_STORAGE_KEY } from '@Domain/preferences';
import { LOGIN_LAYOUT_APPEARANCES } from './constants/loginLayout.appearances';
import { LOGIN_LAYOUT_TEST_IDS } from './constants/loginLayout.testIds';
import { LoginLayoutVariant } from './constants/loginLayout.variant';
import { LOGIN_LAYOUT_COPY } from './constants/loginLayout.copy';
import { getLoginLayoutClassName } from './helpers/loginLayout.getClassName';
import { renderLoginLayoutAppearanceOption } from './renderer/loginLayout.appearanceOption';
import type { LoginLayoutComponent } from './types/loginLayout.props';

import { useAuthViewportFit } from '../../hooks/useAuthViewportFit';
import { LoginLayoutVisual } from './components/visual';
import { getLoginLayoutHeading } from './helpers/loginLayout.getHeading';
import { LoginLayoutFooter } from './components/footer';

// One stylesheet owns the entire auth layout. Legacy variants are intentionally not loaded.
import './styles/optionOneAuth.css';

/** Shared, card-free shell for sign-in and every recovery step. */
export const LoginLayout: LoginLayoutComponent = ({
  children,
  title,
  mobileTitle,
  description,
  mobileDescription,
  minimalMobile = false,
  variant = LoginLayoutVariant.LOGIN,
  eyebrow = LOGIN_LAYOUT_COPY.DEFAULT_EYEBROW,
}) => {
  const { preference, setPreference } = useWorkspaceTheme(SCHOOL_THEME_STORAGE_KEY);
  const { contentRef, toolbarRef, mainRef, isMobile } = useAuthViewportFit();

  return (
    <Box
      as={BOX_TAG.MAIN}
      className={getLoginLayoutClassName(variant, minimalMobile)}
      data-testid={LOGIN_LAYOUT_TEST_IDS.ROOT}
    >
      <Box as={BOX_TAG.DIV} className="rdr-login-v2__frame">
        <LoginLayoutVisual />
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
              aria-label={__('Back to RDR School homepage')}
            >
              <ArrowLeft size={19} />
              <Box as={BOX_TAG.SPAN}>{__('Back to school')}</Box>
            </AppLink>
            <Box
              as={BOX_TAG.DIV}
              className="rdr-login-v2__appearance"
              role="group"
              aria-label={__('Appearance')}
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
                {getLoginLayoutHeading(isMobile, title, mobileTitle)}
              </Box>
              <Box as={BOX_TAG.P} className="rdr-login-v2__description-desktop">
                {description}
              </Box>
              <Box as={BOX_TAG.P} className="rdr-login-v2__description-mobile">
                {mobileDescription ?? description}
              </Box>
            </Box>
            {children}
            <LoginLayoutFooter />
          </Box>
        </Box>
      </Box>
    </Box>
  );
};
