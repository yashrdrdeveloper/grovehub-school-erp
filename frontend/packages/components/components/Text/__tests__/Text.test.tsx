import { render, screen } from '@testing-library/react';
import { Text } from '../Text';
describe('Text', () => {
  it('renders accessible semantic text using the shared typography tokens', () => {
    render(
      <Text as="small" tone="muted" size="sm" dataTestId="school-helper">
        Helper
      </Text>,
    );
    const el = screen.getByTestId('school-helper');
    expect(el.tagName).toBe('SMALL');
    expect(el.textContent).toBe('Helper');
    expect(el.className).toContain('gh-typography--size-sm');
  });
});
