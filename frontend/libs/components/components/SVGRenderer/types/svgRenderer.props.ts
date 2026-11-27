import type { ComponentType, SVGProps } from 'react';
export type SVGRendererSource =
  | { readonly component: ComponentType<SVGProps<SVGSVGElement>>; readonly src?: never }
  | { readonly src: string; readonly component?: never };
export type SVGRendererLabel =
  | { readonly decorative: true; readonly alt?: never }
  | { readonly decorative?: false; readonly alt: string };
export type SVGRendererProps = SVGRendererSource &
  SVGRendererLabel & {
    readonly className?: string;
    readonly dataTestId?: string;
  };
