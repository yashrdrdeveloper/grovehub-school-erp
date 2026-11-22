import type { ReactNode } from 'react';

import { LoaderCircle } from '@Components/icons';

import type { RenderButtonContent } from './types/button.content';

/** Standard loading treatment is rendered once by the primitive, not by each feature. */
export const renderButtonContent: RenderButtonContent = ({
  loading,
  loadingText,
  children,
}) => {
  if (!loading) {
    return children;
  }

  return (
    <>
      <LoaderCircle aria-hidden="true" className="gh-button__spinner" size={16} />
      <span>{loadingText}</span>
    </>
  );
};
