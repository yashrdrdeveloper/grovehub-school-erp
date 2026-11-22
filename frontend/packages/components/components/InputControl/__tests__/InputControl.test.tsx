import { render, screen } from '@testing-library/react';
import { InputControl } from '../InputControl';
describe('InputControl', () => {
  it('forwards a stable selector to the native input and applies control size', () => {
    render(<InputControl aria-label="Search students" dataTestId="student-search" controlSize="sm"/>);
    const input = screen.getByTestId('student-search');
    expect(input.tagName).toBe('INPUT');
    expect(input.className).toContain('gh-input-control--sm');
  });
  it('keeps existing auth layouts when explicitly unstyled', () => {
    render(<InputControl aria-label="Email" dataTestId="login-email" unstyled />);
    expect(screen.getByTestId('login-email').className).not.toContain('gh-input-control');
  });
});
