import { fireEvent, render, screen } from '@testing-library/react';

import { photographs } from '../../../../../constants/school.content';
import { PanelPhoto } from '../PanelPhoto';

describe('PanelPhoto', () => {
  it('retains descriptive alternatives and retries a missing school photograph once', () => {
    render(<PanelPhoto index={0} alt="School achievements" />);

    const photo = screen.getByRole('img', { name: 'School achievements' });
    expect(photo.getAttribute('src')).toBe(photographs[0].src);

    fireEvent.error(photo);
    expect(photo.getAttribute('src')).toBe(photographs[0].originalSrc);

    fireEvent.error(photo);
    expect(photo.hasAttribute('hidden')).toBe(true);
  });
});
