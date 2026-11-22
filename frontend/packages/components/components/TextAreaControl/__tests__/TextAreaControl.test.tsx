import { render, screen } from '@testing-library/react';
import { TextAreaControl } from '../TextAreaControl';

describe('TextAreaControl', () => {
  it('forwards native semantics and a deterministic test ID', () => {
    render(<TextAreaControl aria-label="Notes" dataTestId="notes-field" defaultValue="A note"/>);
    expect(screen.getByRole('textbox',{name:'Notes'})).toBe(screen.getByTestId('notes-field'));
  });
});
