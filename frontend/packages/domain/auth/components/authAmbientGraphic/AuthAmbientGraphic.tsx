import { AcademicPathwaysSvg } from '@Assets/auth';
import { SVGRenderer } from '@Components/components/SVGRenderer';

/** Feature-specific frame; SVG shape is owned by @Assets. */
export const AuthAmbientGraphic = () => (
  <div className="rdr-auth-art" data-testid="auth-ambient-art" aria-hidden="true">
    <SVGRenderer component={AcademicPathwaysSvg} decorative
      className="rdr-auth-art__svg" dataTestId="auth-ambient-svg"/>
    <div className="rdr-auth-art__caption">
      <span>LEARNING</span><span className="rdr-auth-art__sep"/> <span>PEOPLE</span><span className="rdr-auth-art__sep"/> <span>PROGRESS</span>
    </div>
  </div>
);
