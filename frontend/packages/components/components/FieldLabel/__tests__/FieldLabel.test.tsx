import { render, screen } from '@testing-library/react';
import { FieldLabel } from '../FieldLabel';

describe('FieldLabel', () => {
  it('forwards the htmlFor and native test ID', () => {
    render(<><FieldLabel htmlFor="guardian" dataTestId="guardian-label">Guardian</FieldLabel><input id="guardian"/></>);
    expect(screen.getByText('Guardian').getAttribute('for')).toBe('guardian');
    expect(screen.getByTestId('guardian-label').tagName).toBe('LABEL');
  });
});
