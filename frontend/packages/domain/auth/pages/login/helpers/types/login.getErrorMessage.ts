import type { LoginState } from '../../types';

export type GetLoginErrorMessage = (failure: unknown) => LoginState['errorMessage'];
