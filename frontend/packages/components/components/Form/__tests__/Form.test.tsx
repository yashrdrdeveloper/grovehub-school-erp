import { fireEvent, render, screen } from '@testing-library/react';
import { Form } from '../Form';
describe('Form', () => {
  it('forwards its automation selector and native submit handler', () => {
    const handleSubmit = jest.fn(event => event.preventDefault());
    render(<Form dataTestId="school-enquiry" layout="stack" spacing="sm" onSubmit={handleSubmit}>
      <button type="submit">Send</button>
    </Form>);
    const form = screen.getByTestId('school-enquiry');
    expect(form.tagName).toBe('FORM');
    expect(form.className).toContain('gh-form--gap-sm');
    fireEvent.submit(form);
    expect(handleSubmit).toHaveBeenCalledTimes(1);
  });
});
