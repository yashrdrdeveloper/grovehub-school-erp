import { render, screen } from '@testing-library/react';

import { AuthAmbientGraphic } from '../AuthAmbientGraphic';
import { AUTH_AMBIENT_TEST_IDS } from '../constants/authAmbientGraphic.testIds';

describe('AuthAmbientGraphic', () => {
  it('renders one decorative school pathway and its accompanying caption', () => {
    render(<AuthAmbientGraphic />);

    const art = screen.getByTestId(AUTH_AMBIENT_TEST_IDS.ART);
    const svg = screen.getByTestId(AUTH_AMBIENT_TEST_IDS.SVG);

    expect(art.getAttribute('aria-hidden')).toBe('true');
    expect(svg.tagName.toLowerCase()).toBe('svg');
    expect(art.textContent).toContain('LEARNING');
    expect(art.textContent).toContain('PEOPLE');
    expect(art.textContent).toContain('PROGRESS');
  });
});
