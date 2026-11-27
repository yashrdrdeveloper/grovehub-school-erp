import { cn } from '@Base/utilities/cn';
import { RdrSchoolMonogramSvg } from '@Assets/brand';
import { Box, BOX_TAG } from '@Components/components/Box';
import { SVGRenderer } from '@Components/components/SVGRenderer';

import { SCHOOL_MONOGRAM_CLASSES } from './constants/schoolMonogram.classes';
import type { SchoolMonogramProps } from './types/schoolMonogram.props';

/** School-owned, theme-aware decorative brand artwork. */
export const SchoolMonogram = ({ className, dataTestId }: SchoolMonogramProps) => (
  <Box
    as={BOX_TAG.SPAN}
    className={cn(SCHOOL_MONOGRAM_CLASSES.ROOT, className)}
    dataTestId={dataTestId}
    aria-hidden="true"
  >
    <SVGRenderer
      component={RdrSchoolMonogramSvg}
      decorative
      className={SCHOOL_MONOGRAM_CLASSES.SYMBOL}
    />
  </Box>
);
