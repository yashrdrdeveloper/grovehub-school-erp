import type { AuthEndpoint } from '../../constants/auth.endpoints';
import type { AuthOperation } from '../../constants/auth.operations';

export type GetAuthEndpoint = (operation: AuthOperation) => AuthEndpoint;
