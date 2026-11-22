import { useState, type FormEvent } from 'react';
import { ArrowLeft, MailPlus, ShieldCheck, ShieldOff, UserRoundX } from '@Components/icons';
import { Link, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@Components/hooks/query';

import { accountsApi } from '@Business/accounts';
import { Button, Card, Input, PageHeading, StatusPill } from '@Components';

export const AccountDetailsPage = () => {
  const { accountId = '' } = useParams();
  const queryClient = useQueryClient();
  const [username, setUsername] = useState<string | null>(null);
  const [message, setMessage] = useState('');

  const account = useQuery({
    queryKey: ['account', accountId],
    queryFn: () => accountsApi.get(accountId),
    enabled: Boolean(accountId),
  });

  const refresh = () => {
    void account.refetch();
    void queryClient.invalidateQueries({ queryKey: ['accounts'] });
  };

  const update = useMutation({
    mutationFn: (next: string | null) => accountsApi.update(accountId, next),
    onSuccess: () => {
      setMessage('Account updated.');
      refresh();
    },
  });

  const lifecycle = useMutation({
    mutationFn: (action: 'suspend' | 'activate' | 'deactivate') => {
      const actions = {
        suspend: accountsApi.suspend,
        activate: accountsApi.activate,
        deactivate: accountsApi.deactivate,
      };
      return actions[action](accountId);
    },
    onSuccess: () => {
      setMessage('Account status updated.');
      refresh();
    },
  });

  const resend = useMutation({
    mutationFn: () => accountsApi.resendInvitation(accountId),
    onSuccess: () => setMessage('Invitation queued for delivery.'),
  });

  const save = (event: FormEvent) => {
    event.preventDefault();
    update.mutate(username ?? account.data?.username ?? null);
  };

  return (
    <div className="page">
      <Link className="back-link back-link--left" to="/accounts"><ArrowLeft size={16} /> All accounts</Link>
      <PageHeading eyebrow="Identity management" title="Account details" description="Manage the linked school account and account lifecycle." />
      {message ? <div className="notice">{message}</div> : null}
      {account.isLoading ? <Card>Loading account…</Card> : account.error ? (
        <Card><p className="form-error">{account.error.message}</p></Card>
      ) : account.data ? (
        <div className="settings-grid">
          <Card>
            <div className="section-heading"><h3>Account profile</h3><StatusPill value={account.data.status} /></div>
            <dl className="details-list">
              <div><dt>Account ID</dt><dd className="mono">{account.data.accountId}</dd></div>
              <div><dt>Foundation Person GID</dt><dd className="mono">{account.data.personGid}</dd></div>
              <div><dt>Email</dt><dd>{account.data.email}</dd></div>
              <div><dt>Version</dt><dd>{account.data.version}</dd></div>
            </dl>
            <div className="sub-divider" />
            <form onSubmit={save} className="auth-form">
              <Input
                id="accountUsername"
                label="Username"
                value={username ?? account.data.username ?? ''}
                onChange={(event) => setUsername(event.target.value)}
                minLength={3}
                maxLength={40}
                hint="3–40 letters, numbers, dots, underscores or dashes."
              />
              {update.error ? <p className="form-error">{update.error.message}</p> : null}
              <Button type="submit" disabled={update.isPending}>Save username</Button>
            </form>
          </Card>
          <Card>
            <div className="section-heading"><h3>Account actions</h3></div>
            <div className="action-stack">
              <Button variant="secondary" disabled={resend.isPending || account.data.status !== 'PENDING_SETUP'} onClick={() => resend.mutate()}><MailPlus size={16} /> Resend invitation</Button>
              <Button variant="secondary" disabled={lifecycle.isPending || account.data.status === 'ACTIVE'} onClick={() => lifecycle.mutate('activate')}><ShieldCheck size={16} /> Activate account</Button>
              <Button variant="secondary" disabled={lifecycle.isPending || account.data.status !== 'ACTIVE'} onClick={() => lifecycle.mutate('suspend')}><ShieldOff size={16} /> Suspend account</Button>
              <Button variant="danger" disabled={lifecycle.isPending || account.data.status === 'DEACTIVATED'} onClick={() => {
                if (window.confirm('Deactivate this account? This will revoke its sessions.')) {
                  lifecycle.mutate('deactivate');
                }
              }}><UserRoundX size={16} /> Deactivate account</Button>
            </div>
            {lifecycle.error || resend.error ? <p className="form-error">{(lifecycle.error ?? resend.error)?.message}</p> : null}
          </Card>
        </div>
      ) : null}
    </div>
  );
};
