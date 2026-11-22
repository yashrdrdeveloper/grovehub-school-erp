import { ArrowLeft, LockKeyhole } from '@Components/icons';
import { Link } from 'react-router-dom';
import { APP_ROUTE } from '@Business/app/routes';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme } from '@Components/hooks/theme';
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
    <main className={getLoginLayoutClassName(variant, minimalMobile)} data-testid={LOGIN_LAYOUT_TEST_IDS.ROOT}>
      <div className="rdr-login-v2__frame">
        <section className="rdr-login-v2__visual" aria-label="RDR School branded gradient background">
          <div className="rdr-login-v2__mobile-header">
            <div className="rdr-login-v2__mobile-identity" aria-label="RDR School">
              <SchoolMonogram/>
              <span className="rdr-login-v2__mobile-school-name">RDR School<small>Your school workspace</small></span>
            </div>
          </div>
          <div className="rdr-login-v2__visual-brand">
            <SchoolMonogram/>
            <span><strong>RDR School</strong><small>Rao Deena Ram Vidya Vihar</small></span>
          </div>
          <AuthAmbientGraphic/>
          <div className="rdr-login-v2__visual-copy">
            <span className="rdr-login-v2__visual-kicker"><span/> THE RDR SCHOOL WORKSPACE</span>
            <h2>For every school day.<br/><em>And everyone behind it.</em></h2>
            <p>A connected space for the people who make learning happen.</p>
            <div className="rdr-login-v2__visual-footer"><span>LEARN</span><i aria-hidden="true"/><span>ORGANISE</span><i aria-hidden="true"/><span>GROW</span></div>
          </div>
        </section>
        <section ref={contentRef} className="rdr-login-v2__content" data-testid={LOGIN_LAYOUT_TEST_IDS.CONTENT} aria-labelledby="rdr-auth-title">
          <div ref={toolbarRef} className="rdr-login-v2__toolbar" data-testid={LOGIN_LAYOUT_TEST_IDS.TOOLBAR}>
            <Link to={APP_ROUTE.HOME} className="rdr-login-v2__back" aria-label="Back to RDR School homepage">
              <ArrowLeft size={19}/><span>Back to school</span>
            </Link>
            <div className="rdr-login-v2__appearance" role="group" aria-label="Appearance">
              {LOGIN_LAYOUT_APPEARANCES.map(option => renderLoginLayoutAppearanceOption({ ...option, preference, setPreference }))}
            </div>
          </div>
          <div ref={mainRef} className="rdr-login-v2__main" data-testid={LOGIN_LAYOUT_TEST_IDS.MAIN}>
            <div className="rdr-login-v2__heading">
              <span className="rdr-login-v2__eyebrow">{eyebrow}</span>
              <h1 id="rdr-auth-title" data-testid={LOGIN_LAYOUT_TEST_IDS.TITLE}>{title}</h1>
              <p className="rdr-login-v2__description-desktop">{description}</p>
              <p className="rdr-login-v2__description-mobile">{mobileDescription ?? description}</p>
            </div>
            {children}
            <div className="rdr-login-v2__closing" aria-label="Secure RDR School workspace">
              <LockKeyhole size={13} strokeWidth={1.8} aria-hidden="true"/>
              <span>A secure workspace for RDR School</span>
            </div>
            <div className="rdr-o1-auth-footer"><span>© RDR School</span><span>Haluhera · Rewari</span></div>
          </div>
        </section>
      </div>
    </main>
  );
};
