import { getFormFieldDescriptionId } from '../formField.getDescriptionId';

describe('getFormFieldDescriptionId', () => {
  it('creates an accessible description ID when error or hint exists', () => {
    expect(getFormFieldDescriptionId('email', 'Invalid email')).toBe('email-description');
    expect(getFormFieldDescriptionId('email', undefined, 'School email only'))
      .toBe('email-description');
  });

  it('omits an ID when no description is rendered', () => {
    expect(getFormFieldDescriptionId('email')).toBeUndefined();
  });
});
