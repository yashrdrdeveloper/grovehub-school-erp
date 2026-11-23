import { render, screen } from '@testing-library/react';

import { LoginLayoutVisual } from '../LoginLayoutVisual';

describe('LoginLayoutVisual', () => {
  it('preserves RDR branding and the decorative visual section', () => {
    render(<LoginLayoutVisual />);

    expect(screen.getByLabelText('RDR School branded gradient background')).toBeTruthy();
    expect(screen.getByText('Rao Deena Ram Vidya Vihar')).toBeTruthy();
  });
});
