import { getFormFieldControlTestId } from '../formField.getControlTestId';

describe('getFormFieldControlTestId', () => {
  it('uses explicit test ID before element ID', () => {
    expect(
      getFormFieldControlTestId({
        label: 'Email address',
        id: 'email',
        dataTestId: 'account-email',
      }),
    ).toBe('account-email');
  });

  it('uses the explicit element ID next', () => {
    expect(getFormFieldControlTestId({ label: 'Email address', id: 'school-email' })).toBe(
      'school-email',
    );
  });

  it('normalizes labels using Lodash', () => {
    expect(getFormFieldControlTestId({ label: "Parent's E-mail Address!" })).toBe(
      'parents-e-mail-address',
    );
  });
});
