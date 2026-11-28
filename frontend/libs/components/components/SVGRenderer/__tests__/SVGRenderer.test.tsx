import type { SVGProps } from 'react';
import { render, screen } from '@testing-library/react';
import { SVGRenderer } from '../SVGRenderer';
const Sample = (props: SVGProps<SVGSVGElement>) => (
  <svg {...props} viewBox="0 0 10 10">
    <circle cx="5" cy="5" r="5" />
  </svg>
);
describe('SVGRenderer', () => {
  it('renders an accessible inline SVG with a native test ID', () => {
    render(<SVGRenderer component={Sample} alt="Sample symbol" dataTestId="sample-symbol" />);
    expect(screen.getByRole('img', { name: 'Sample symbol' }).getAttribute('data-testid')).toBe(
      'sample-symbol',
    );
  });
  it('renders file SVGs as decorative images when requested', () => {
    render(<SVGRenderer src="/sample.svg" decorative dataTestId="sample-art" />);
    expect(screen.getByTestId('sample-art').getAttribute('alt')).toBe('');
  });
});
