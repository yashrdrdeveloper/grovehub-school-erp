import { render, screen } from '@testing-library/react';

import { renderButtonContent } from '../button.content';

describe('renderButtonContent', () => {
  it('renders standard content without a spinner', () => {
    const { container } = render(
      <>{renderButtonContent({ loading: false, loadingText: 'Working', children: 'Continue' })}</>,
    );

    expect(screen.getByText('Continue')).toBeTruthy();
    expect(container.querySelector('.gh-button__spinner')).toBeNull();
  });

  it('renders an accessible loading label with a decorative spinner', () => {
    const { container } = render(
      <>{renderButtonContent({ loading: true, loadingText: 'Working', children: 'Continue' })}</>,
    );

    expect(screen.getByText('Working')).toBeTruthy();
    expect(container.querySelector('.gh-button__spinner')).toBeTruthy();
  });
});
