import { render, screen } from '@testing-library/react';

import { LoginLayoutFooter } from '../LoginLayoutFooter';

describe('LoginLayoutFooter', () => {
  it('shows the same secure workspace and school location copy', () => {
    render(<LoginLayoutFooter />);

    expect(screen.getByText('A secure workspace for RDR School')).toBeTruthy();
    expect(screen.getByText('Haluhera · Rewari')).toBeTruthy();
  });
});
