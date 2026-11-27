import { apiRequest } from '@Business/api';
import { HTTP_METHOD } from '@Business/api/constants/api.method';

import { INTEGRATION_ENDPOINT } from './constants/integrations.endpoints';
import { getIntegrationProviderPath } from './helpers/integrations.getProviderPath';
import { getIntegrationResourcePath } from './helpers/integrations.getResourcePath';

/** Integration transport facade, independent of browser routes or UI. */
export const integrationApi = {
  providers: () => apiRequest<unknown[]>(INTEGRATION_ENDPOINT.PUBLIC_PROVIDERS),

  provider: (provider: string) => apiRequest<unknown>(getIntegrationProviderPath(provider)),

  create: (input: unknown) =>
    apiRequest<unknown>(INTEGRATION_ENDPOINT.PROTECTED_CONNECTIONS, {
      method: HTTP_METHOD.POST,
      body: input,
    }),

  list: () => apiRequest<unknown[]>(INTEGRATION_ENDPOINT.PROTECTED_CONNECTIONS),

  get: (id: string) => apiRequest<unknown>(getIntegrationResourcePath(id)),

  disconnect: (id: string) =>
    apiRequest<unknown>(getIntegrationResourcePath(id), { method: HTTP_METHOD.DELETE }),

  beginOAuth: (id: string) =>
    apiRequest<unknown>(`${getIntegrationResourcePath(id)}/oauth/start`, {
      method: HTTP_METHOD.POST,
    }),

  runTest: (id: string, direction: string) =>
    apiRequest<unknown>(`${getIntegrationResourcePath(id)}/test`, {
      method: HTTP_METHOD.POST,
      body: { direction },
    }),

  tests: (id: string) => apiRequest<unknown[]>(`${getIntegrationResourcePath(id)}/tests`),

  watch: (id: string) =>
    apiRequest<unknown>(`${getIntegrationResourcePath(id)}/watch`, {
      method: HTTP_METHOD.POST,
    }),

  executeAction: (id: string, action: string, input: unknown) =>
    apiRequest<unknown>(`${getIntegrationResourcePath(id)}/actions/${encodeURIComponent(action)}`, {
      method: HTTP_METHOD.POST,
      body: { input },
    }),
};
