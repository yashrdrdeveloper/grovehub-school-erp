import { render, screen } from '@testing-library/react';
import { Text, TEXT_TAG, TEXT_TONE, TEXT_SIZE } from '../index';
describe('Text', () => {
  it('renders accessible semantic text using the shared typography tokens', () => {
    render(
      <Text
        as={TEXT_TAG.SMALL}
        tone={TEXT_TONE.MUTED}
        size={TEXT_SIZE.SMALL}
        dataTestId="school-helper"
      >
        Helper
      </Text>,
    );
    const el = screen.getByTestId('school-helper');
    expect(el.tagName).toBe('SMALL');
    expect(el.textContent).toBe('Helper');
    expect(el.className).toContain('gh-typography--size-sm');
  });
});
