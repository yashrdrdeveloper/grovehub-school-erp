import { cn } from '@Base/utilities/cn';
import { RdrSchoolMonogramSvg } from '@Assets/brand';
import { SVGRenderer } from './SVGRenderer';

export interface SchoolMonogramProps {
  readonly className?: string;
  readonly dataTestId?: string;
}

/** RDR School's approved merged monogram; artwork remains an inline theme-aware asset. */
export const SchoolMonogram = ({ className, dataTestId }: SchoolMonogramProps) => (
  <span className={cn('school-monogram', className)} data-testid={dataTestId} aria-hidden="true">
    <SVGRenderer component={RdrSchoolMonogramSvg} decorative className="school-monogram__symbol"/>
  </span>
);
