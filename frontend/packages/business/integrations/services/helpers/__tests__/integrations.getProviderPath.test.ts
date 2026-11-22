import { getIntegrationProviderPath } from '../integrations.getProviderPath';

describe('getIntegrationProviderPath', () => {
  it('encodes arbitrary provider identifiers on a public API path', () => {
    expect(getIntegrationProviderPath('gmail/oauth')).toBe('/integrations/public/providers/gmail%2Foauth');
  });
});
