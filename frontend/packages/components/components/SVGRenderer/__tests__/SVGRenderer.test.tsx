import type { SVGProps } from 'react';
import { render, screen } from '@testing-library/react';
import { SVGRenderer } from '../SVGRenderer';
import { SchoolMonogram } from '../../SchoolMonogram';
const Sample = (props: SVGProps<SVGSVGElement>) => (
  <svg {...props} viewBox="0 0 10 10">
    <circle cx="5" cy="5" r="5" />
  </svg>
);
describe('SVGRenderer', () => {
  it('renders an accessible inline SVG with a native test ID', () => {
    render(<SVGRenderer component={Sample} alt="School symbol" dataTestId="school-symbol" />);
    expect(screen.getByRole('img', { name: 'School symbol' }).getAttribute('data-testid')).toBe(
      'school-symbol',
    );
  });
  it('renders file SVGs as decorative images when requested', () => {
    render(<SVGRenderer src="/school.svg" decorative dataTestId="school-art" />);
    expect(screen.getByTestId('school-art').getAttribute('alt')).toBe('');
  });
  it('preserves the connected RDR monogram paths when rendered via the asset package', () => {
    render(<SchoolMonogram dataTestId="school-logo" />);
    const logo = screen.getByTestId('school-logo');
    expect(logo.querySelectorAll('svg.school-monogram__symbol path')).toHaveLength(3);
  });
});
