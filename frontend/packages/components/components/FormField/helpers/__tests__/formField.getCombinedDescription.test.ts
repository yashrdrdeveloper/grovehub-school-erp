import { getFormFieldCombinedDescription } from '../formField.getCombinedDescription';

describe('getFormFieldCombinedDescription', () => {
  it('joins the caller and generated descriptions in their original order', () => {
    expect(getFormFieldCombinedDescription('custom-instructions', 'school-email-description'))
      .toBe('custom-instructions school-email-description');
  });

  it('removes empty values without introducing spaces', () => {
    expect(getFormFieldCombinedDescription('', 'required-message')).toBe('required-message');
    expect(getFormFieldCombinedDescription(undefined, undefined)).toBeUndefined();
  });
});
