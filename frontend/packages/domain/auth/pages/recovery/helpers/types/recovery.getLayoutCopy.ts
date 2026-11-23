export type GetRecoveryLayoutCopy = (params: {
  readonly invitation: boolean;
  readonly formVisible: boolean;
}) => { readonly mobileTitle: string; readonly description: string; readonly mobileDescription: string };
