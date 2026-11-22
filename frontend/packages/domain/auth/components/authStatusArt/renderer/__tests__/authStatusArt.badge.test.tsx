import { render, screen } from '@testing-library/react';
import { renderAuthStatusBadge } from '../authStatusArt.badge';

describe('renderAuthStatusBadge', () => {
  it('renders a success badge with semantic color tokens', () => {
    const { container } = render(<>{renderAuthStatusBadge('success')}</>);
    expect(container.querySelector('svg')?.getAttribute('fill')).toBe('var(--gh-color-success)');
  });
  it('shows an expired marker but no badge during checking', () => {
    render(<>{renderAuthStatusBadge('expired')}</>);
    expect(screen.getByText('×')).toBeTruthy();
    const { container } = render(<>{renderAuthStatusBadge('checking')}</>);
    expect(container.querySelector('svg')).toBeNull();
  });
});
