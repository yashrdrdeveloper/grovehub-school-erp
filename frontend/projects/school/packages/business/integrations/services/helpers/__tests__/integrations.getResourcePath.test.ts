import { getIntegrationResourcePath } from '../integrations.getResourcePath';

describe('getIntegrationResourcePath', () => {
  it('escapes connection identifiers without altering the protected route', () => {
    expect(getIntegrationResourcePath('tenant/connection')).toBe(
      '/integrations/protected/connections/tenant%2Fconnection',
    );
  });
});
