import { cn } from '@Base/utilities';

export interface SchoolMonogramProps {
  className?: string;
}

/** Shared school identity mark. The letters intentionally overlap as a single monogram. */
export const SchoolMonogram = ({ className }: SchoolMonogramProps) => (
  <span className={cn('school-monogram', className)} aria-hidden="true">
    <span className="school-monogram__letters">RDR</span>
  </span>
);
