import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Cable, ExternalLink, RefreshCw, Unplug } from 'lucide-react';
import { useQuery, useQueryClient } from '@tanstack/react-query';

import { integrationApi } from '@Business/integrations';
import { Button, Card, EmptyState, PageHeading, StatusPill } from '@Components';
import { useAuth } from '@Domain/auth';

interface Connection {
  id: string;
  provider: string;
  ownerGid: string;
  status: string;
}

interface Provider {
  id: string;
  name?: string;
}

export const IntegrationsPage = () => {
  const { account } = useAuth();
  const queryClient = useQueryClient();
  const [feedback, setFeedback] = useState('');
  const connections = useQuery({
    queryKey: ['integrations'],
    queryFn: () => integrationApi.list() as Promise<Connection[]>,
  });
  const providers = useQuery({
    queryKey: ['integrationProviders'],
    queryFn: () => integrationApi.providers() as Promise<Provider[]>,
  });
  const [connecting, setConnecting] = useState(false);

  const connect = async (provider: string) => {
    setConnecting(true);
    try {
      const response = await integrationApi.create({ provider, ownerGid: account?.personGid }) as Connection;
      const details = await integrationApi.beginOAuth(response.id) as { authorizationUrl?: string };
      if (details.authorizationUrl?.startsWith('https://')) {
        window.location.assign(details.authorizationUrl);
      } else {
        setFeedback('Connection created; provider authorization URL was not returned.');
        void connections.refetch();
      }
    } catch (error) {
      setFeedback(error instanceof Error ? error.message : 'Could not connect provider.');
    } finally {
      setConnecting(false);
    }
  };

  const disconnect = async (id: string) => {
    try {
      await integrationApi.disconnect(id);
      setFeedback('Connection disconnected.');
      void queryClient.invalidateQueries({ queryKey: ['integrations'] });
    } catch (error) {
      setFeedback(error instanceof Error ? error.message : 'Disconnect failed.');
    }
  };

  return (
    <div className="page">
      <PageHeading eyebrow="Workspace connections" title="Integrations" description="Connect approved mail and provider services to your school workspace." />
      {feedback ? <div className="notice">{feedback}</div> : null}
      <div className="settings-grid">
        <Card>
          <div className="section-heading"><h3>Connected services</h3><Button variant="ghost" onClick={() => void connections.refetch()}><RefreshCw size={16} /></Button></div>
          {connections.isLoading ? <p className="muted">Loading connections…</p> : connections.error ? (
            <p className="form-error">{connections.error.message}</p>
          ) : !connections.data?.length ? (
            <EmptyState title="No connections yet" description="Select an available provider to get started." />
          ) : connections.data.map((connection) => (
            <div className="session-item" key={connection.id}>
              <Cable size={20} />
              <div><strong><Link to={`/integrations/${connection.id}`}>{connection.provider}</Link></strong><small><StatusPill value={connection.status} /></small></div>
              <Button variant="ghost" onClick={() => void disconnect(connection.id)}><Unplug size={17} /></Button>
            </div>
          ))}
        </Card>
        <Card>
          <div className="section-heading"><h3>Available providers</h3><span className="muted">Choose a connection</span></div>
          {providers.isLoading ? <p className="muted">Loading providers…</p> : providers.error ? (
            <p className="form-error">{providers.error.message}</p>
          ) : !providers.data?.length ? (
            <EmptyState title="No providers configured" description="Integration providers must be enabled on the backend." />
          ) : providers.data.map((provider) => (
            <div className="session-item" key={provider.id}>
              <Cable size={20} />
              <div><strong>{provider.name ?? provider.id}</strong><small>OAuth authorization</small></div>
              <Button variant="secondary" disabled={connecting} onClick={() => void connect(provider.id)}>Connect <ExternalLink size={15} /></Button>
            </div>
          ))}
        </Card>
      </div>
    </div>
  );
};
