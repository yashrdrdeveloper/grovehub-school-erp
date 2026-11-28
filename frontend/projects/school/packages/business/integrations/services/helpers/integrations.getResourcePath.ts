import { INTEGRATION_ENDPOINT } from '../constants/integrations.endpoints';

import type { GetIntegrationResourcePath } from './types/integrations.getResourcePath';

export const getIntegrationResourcePath: GetIntegrationResourcePath = (id) =>
  `${INTEGRATION_ENDPOINT.PROTECTED_CONNECTIONS}/${encodeURIComponent(id)}`;
