import type { ReactElement } from 'react';
import type { SVGRendererProps } from './types/svgRenderer.props';
/** Native accessible leaf for file-backed SVGs and CSS-inheriting inline SVGs. */
export const SVGRenderer = (props: SVGRendererProps): ReactElement => {
  const { className, dataTestId, decorative = false } = props;
  if (props.component) {
    const Graphic = props.component;
    if (decorative) return <Graphic className={className} aria-hidden="true" focusable="false" data-testid={dataTestId}/>;
    return <Graphic className={className} role="img" aria-label={props.alt} focusable="false" data-testid={dataTestId}/>;
  }
  if (decorative) return <img src={props.src} className={className} alt="" aria-hidden="true" data-testid={dataTestId}/>;
  return <img src={props.src} className={className} alt={props.alt} data-testid={dataTestId}/>;
};
