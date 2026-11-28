import { BOX_TAG } from '@Components/components/Box';
import { createRef } from 'react';
import { render, screen } from '@testing-library/react';
import { Box } from '../Box';

describe('Box', () => {
  it('forwards a test ID and ref without changing native section semantics', () => {
    const ref = createRef<HTMLElement>();
    render(
      <Box
        as={BOX_TAG.SECTION}
        role="region"
        aria-label="School sections"
        ref={ref}
        dataTestId="school-section"
      >
        <span>Content</span>
      </Box>,
    );
    expect(screen.getByRole('region', { name: 'School sections' })).toBe(ref.current);
    expect(screen.getByTestId('school-section').tagName).toBe('SECTION');
  });
});
