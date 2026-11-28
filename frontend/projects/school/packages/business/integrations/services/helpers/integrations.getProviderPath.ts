import { INTEGRATION_ENDPOINT } from '../constants/integrations.endpoints';

import type { GetIntegrationProviderPath } from './types/integrations.getProviderPath';

export const getIntegrationProviderPath: GetIntegrationProviderPath = (provider) =>
  `${INTEGRATION_ENDPOINT.PUBLIC_PROVIDERS}/${encodeURIComponent(provider)}`;
