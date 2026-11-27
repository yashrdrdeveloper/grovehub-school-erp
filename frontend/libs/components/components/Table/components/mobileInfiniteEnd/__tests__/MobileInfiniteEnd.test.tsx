import { render, screen } from '@testing-library/react';
import { MobileInfiniteEnd } from '../MobileInfiniteEnd';

describe('MobileInfiniteEnd', () => {
  it('renders nothing after the last page', () => {
    const { container } = render(<MobileInfiniteEnd hasNext={false} loading={false} />);

    expect(container.firstChild).toBeNull();
  });

  it('keeps the explicit load-more fallback while further records exist', () => {
    render(<MobileInfiniteEnd hasNext loading={false} />);

    expect(screen.getByRole('button', { name: 'Load more' })).toBeTruthy();
  });
});
