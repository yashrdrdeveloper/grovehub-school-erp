import { fireEvent, render, screen } from '@testing-library/react';
import { Button } from '../Button';
import { BUTTON_VARIANT } from '../constants/button.variant';
import { BUTTON_SIZE } from '../constants/button.size';

describe('Button', () => {
  it('uses exported variant and size constants in the CVA recipe', () => {
    render(
      <Button variant={BUTTON_VARIANT.OUTLINE} size={BUTTON_SIZE.LARGE}>
        Review
      </Button>,
    );
    const button = screen.getByRole('button', { name: 'Review' });
    expect(button.className).toContain('button--outline');
    expect(button.className).toContain('gh-button--lg');
  });
  it('retains native focus and classes without CVA overrides in unstyled mode', () => {
    const onClick = jest.fn();
    render(
      <Button unstyled type="button" className="legacy" dataTestId="visibility" onClick={onClick}>
        Reveal
      </Button>,
    );
    const button = screen.getByRole('button', { name: 'Reveal' });
    expect(button.className).toBe('legacy');
    expect(screen.getByTestId('visibility')).toBe(button);
    fireEvent.click(button);
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
