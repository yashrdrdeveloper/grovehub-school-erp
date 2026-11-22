import { fireEvent, render, screen } from '@testing-library/react';
import { Button } from '../Button';

describe('Button', () => {
  it('retains native focus and classes without CVA overrides in unstyled mode', () => {
    const onClick = jest.fn();
    render(<Button unstyled type="button" className="legacy" dataTestId="visibility" onClick={onClick}>Reveal</Button>);
    const button = screen.getByRole('button',{name:'Reveal'});
    expect(button.className).toBe('legacy');
    expect(screen.getByTestId('visibility')).toBe(button);
    fireEvent.click(button);
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
