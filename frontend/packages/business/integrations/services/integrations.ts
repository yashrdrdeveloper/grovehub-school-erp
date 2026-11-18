import { apiRequest } from '@Business/api';

const ROOT = '/integrations/protected/connections';
const RESOURCE = (id: string) => `${ROOT}/${encodeURIComponent(id)}`;

export const integrationApi = {
  providers: () => apiRequest<unknown[]>('/integrations/public/providers'),
  provider: (provider: string) =>
    apiRequest<unknown>(`/integrations/public/providers/${encodeURIComponent(provider)}`),
  create: (input: unknown) => apiRequest<unknown>(ROOT, { method: 'POST', body: input }),
  list: () => apiRequest<unknown[]>(ROOT),
  get: (id: string) => apiRequest<unknown>(RESOURCE(id)),
  disconnect: (id: string) => apiRequest<unknown>(RESOURCE(id), { method: 'DELETE' }),
  beginOAuth: (id: string) =>
    apiRequest<unknown>(`${RESOURCE(id)}/oauth/start`, { method: 'POST' }),
  runTest: (id: string, direction: string) =>
    apiRequest<unknown>(`${RESOURCE(id)}/test`, { method: 'POST', body: { direction } }),
  tests: (id: string) => apiRequest<unknown[]>(`${RESOURCE(id)}/tests`),
  watch: (id: string) => apiRequest<unknown>(`${RESOURCE(id)}/watch`, { method: 'POST' }),
  executeAction: (id: string, action: string, input: unknown) =>
    apiRequest<unknown>(`${RESOURCE(id)}/actions/${encodeURIComponent(action)}`, { method: 'POST', body: { input } }),
};
