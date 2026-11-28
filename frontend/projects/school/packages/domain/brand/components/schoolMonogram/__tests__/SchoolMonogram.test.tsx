import { render, screen } from '@testing-library/react';
import { SchoolMonogram } from '../SchoolMonogram';

describe('SchoolMonogram', () => {
  it('preserves the decorative span, school artwork and forwarded selector', () => {
    render(<SchoolMonogram className="custom-brand" dataTestId="school-brand" />);
    const root = screen.getByTestId('school-brand');

    expect(root.tagName).toBe('SPAN');
    expect(root.classList.contains('school-monogram')).toBe(true);
    expect(root.classList.contains('custom-brand')).toBe(true);
    expect(root.getAttribute('aria-hidden')).toBe('true');
    expect(root.querySelectorAll('svg.school-monogram__symbol path')).toHaveLength(3);
  });
});
