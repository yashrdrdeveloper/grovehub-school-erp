import { render, screen } from '@testing-library/react';

import { StatusPill } from '../StatusPill';

describe('StatusPill', () => {
  it('converts a backend status into an accessible text label', () => {
    render(<StatusPill value="PENDING_SETUP" dataTestId="account-status" />);

    expect(screen.getByTestId('account-status').textContent).toBe('pending setup');
    expect(screen.getByTestId('account-status').className).toContain('status-pill--notice');
  });
});
