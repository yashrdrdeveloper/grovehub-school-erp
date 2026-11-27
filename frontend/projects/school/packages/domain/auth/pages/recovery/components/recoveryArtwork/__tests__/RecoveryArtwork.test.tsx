import { render, screen } from '@testing-library/react';

import { RecoveryArtwork } from '../RecoveryArtwork';

describe('RecoveryArtwork', () => {
  it('preserves accessible recovery status without the hidden legacy art subtree', () => {
    render(<RecoveryArtwork view="expired" purpose="recovery" />);

    const status = screen.getByRole('status');
    expect(status.getAttribute('aria-live')).toBe('polite');
    expect(screen.getByTestId('auth-status-expired')).toBeTruthy();
    expect(status.querySelectorAll('strong')).toHaveLength(2);
  });
});
