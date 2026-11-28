/** Keeps mobile and desktop recovery messages coherent without JSX mutation. */
export const getRecoveryLayoutCopy = ({
  invitation,
  formVisible,
}: {
  readonly invitation: boolean;
  readonly formVisible: boolean;
}) => {
  const description = __('Choose a strong password to protect your school account.');
  const mobileTitle = (invitation && __('Set password')) || __('Reset password');
  const mobileDescription = (formVisible && description) || '';
  return { mobileTitle, description, mobileDescription };
};
