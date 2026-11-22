import { render, screen } from '@testing-library/react';
import { Card } from '../Card';
describe('Card', () => {
  it('applies a semantic variant and padding from shared design tokens', () => {
    render(<Card dataTestId="school-card" variant="outline" padding="sm">School</Card>);
    const card = screen.getByTestId('school-card');
    expect(card.className).toContain('gh-card--outline');
    expect(card.className).toContain('gh-card--padding-sm');
  });
});
