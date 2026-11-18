import { useState, type FormEvent } from 'react';
import { ArrowLeft, FlaskConical, RadioTower, Send } from 'lucide-react';
import { Link, useParams } from 'react-router-dom';
import { useMutation, useQuery } from '@Components/hooks/query';

import { integrationApi } from '@Business/integrations';
import { Button, Card, PageHeading, StatusPill } from '@Components';

interface Connection {
  id: string;
  provider: string;
  status: string;
  ownerGid: string;
  externalId?: string;
}

interface ConnectionTest {
  id: string;
  direction: string;
  status: string;
  detail: string;
  startedAt: string;
}

export const ConnectionDetailsPage = () => {
  const { connectionId = '' } = useParams();
  const [action, setAction] = useState('mail.send');
  const [payload, setPayload] = useState('{}');
  const [feedback, setFeedback] = useState('');

  const details = useQuery({
    queryKey: ['connection', connectionId],
    queryFn: () => integrationApi.get(connectionId) as Promise<Connection>,
    enabled: Boolean(connectionId),
  });

  const tests = useQuery({
    queryKey: ['connectionTests', connectionId],
    queryFn: () => integrationApi.tests(connectionId) as Promise<ConnectionTest[]>,
    enabled: Boolean(connectionId),
  });

  const runTest = useMutation({
    mutationFn: (direction: string) => integrationApi.runTest(connectionId, direction),
    onSuccess: () => { setFeedback('Test requested.'); void tests.refetch(); },
  });

  const watch = useMutation({
    mutationFn: () => integrationApi.watch(connectionId),
    onSuccess: () => setFeedback('Provider watch updated.'),
  });

  const execute = useMutation({
    mutationFn: () => integrationApi.executeAction(connectionId, action, JSON.parse(payload) as unknown),
    onSuccess: () => setFeedback('Provider action accepted.'),
  });

  const submit = (event: FormEvent) => {
    event.preventDefault();
    try {
      JSON.parse(payload);
      execute.mutate();
    } catch {
      setFeedback('Enter valid JSON action input.');
    }
  };

  return (
    <div className="page">
      <Link to="/integrations" className="back-link back-link--left"><ArrowLeft size={16} /> All integrations</Link>
      <PageHeading eyebrow="Provider connection" title={details.data?.provider ?? 'Connection details'} description="Inspect and operate a configured provider connection." />
      {feedback ? <div className="notice">{feedback}</div> : null}
      {details.error ? <Card><p className="form-error">{details.error.message}</p></Card> : null}
      <div className="settings-grid">
        <Card>
          <div className="section-heading"><h3>Connection</h3>{details.data ? <StatusPill value={details.data.status} /> : null}</div>
          {details.isLoading ? <p className="muted">Loading connection…</p> : (
            <dl className="details-list">
              <div><dt>Connection ID</dt><dd className="mono">{details.data?.id}</dd></div>
              <div><dt>Provider</dt><dd>{details.data?.provider}</dd></div>
              <div><dt>Owner GID</dt><dd className="mono">{details.data?.ownerGid}</dd></div>
            </dl>
          )}
          <div className="sub-divider" />
          <div className="action-stack">
            <Button variant="secondary" loading={runTest.isPending} onClick={() => runTest.mutate('OUTBOUND')}><FlaskConical size={16} /> Run outbound test</Button>
            <Button variant="secondary" loading={watch.isPending} onClick={() => watch.mutate()}><RadioTower size={16} /> Renew watch</Button>
          </div>
          {runTest.error || watch.error ? <p className="form-error">{(runTest.error ?? watch.error)?.message}</p> : null}
        </Card>
        <Card>
          <div className="section-heading"><h3>Provider action</h3></div>
          <p className="muted">Only use actions supported by this provider and your permissions.</p>
          <form className="auth-form" onSubmit={submit}>
            <label className="field"><span className="field__label">Action key</span><input className="field__input" value={action} onChange={(event) => setAction(event.target.value)} required /></label>
            <label className="field"><span className="field__label">Action input (JSON)</span><textarea className="field__input action-json" value={payload} onChange={(event) => setPayload(event.target.value)} required /></label>
            {execute.error ? <p className="form-error">{execute.error.message}</p> : null}
            <Button type="submit" loading={execute.isPending}><Send size={16} /> Execute action</Button>
          </form>
        </Card>
      </div>
      <Card>
        <div className="section-heading"><h3>Recent tests</h3><Button variant="ghost" onClick={() => void tests.refetch()}>Refresh</Button></div>
        {tests.isLoading ? <p className="muted">Loading tests…</p> : tests.error ? <p className="form-error">{tests.error.message}</p> : !tests.data?.length ? (
          <p className="muted">No tests recorded.</p>
        ) : tests.data.map((test) => (
          <div className="session-item" key={test.id}>
            <FlaskConical size={18} /><div><strong>{test.direction}</strong><small>{test.detail || new Date(test.startedAt).toLocaleString()}</small></div><StatusPill value={test.status} />
          </div>
        ))}
      </Card>
    </div>
  );
};
