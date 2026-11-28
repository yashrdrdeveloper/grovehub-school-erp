import { render, screen } from '@testing-library/react';

import { DashboardNotice } from '../DashboardNotice';

describe('DashboardNotice', () => {
  it('discloses preview data honestly instead of inventing live school activity', () => {
    render(<DashboardNotice />);

    const notice = screen.getByRole('status');

    expect(notice.textContent).toContain('Your workspace is taking shape');
    expect(notice.textContent).toContain('This is a preview');
    expect(notice.textContent).toContain('PREVIEW');
  });
});
