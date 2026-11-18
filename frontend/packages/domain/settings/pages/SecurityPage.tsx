import { useState, type FormEvent } from 'react';
import { KeyRound, Lock, MonitorSmartphone, ShieldCheck, Trash2 } from 'lucide-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { authenticationApi } from '@Business/auth';
import { Button, Card, Input, PageHeading } from '@Components';

export const SecurityPage = () => {
  const queryClient = useQueryClient();
  const [newPassword, setNewPassword] = useState('');
  const [oldPassword, setOldPassword] = useState('');
  const [feedback, setFeedback] = useState('');
  const sessions = useQuery({ queryKey: ['sessions'], queryFn: authenticationApi.sessions });
  const reauth = useMutation({
    mutationFn: () => authenticationApi.reauthenticate(oldPassword),
    onSuccess: () => { setOldPassword(''); setFeedback('Identity reverified for sensitive operations.'); },
  });
  const change = useMutation({
    mutationFn: () => authenticationApi.changePassword(newPassword),
    onSuccess: () => {
      setNewPassword('');
      setFeedback('Password changed. Existing sessions have been revoked; sign in again.');
      void queryClient.invalidateQueries({ queryKey: ['sessions'] });
    },
  });
  const revoke = useMutation({
    mutationFn: (sessionId: string) => authenticationApi.revokeSession(sessionId),
    onSuccess: () => void queryClient.invalidateQueries({ queryKey: ['sessions'] }),
  });

  const submitPassword = (event: FormEvent) => {
    event.preventDefault();
    change.mutate();
  };

  return (
    <div className="page">
      <PageHeading eyebrow="Account settings" title="Security & sessions" description="Manage your credentials and signed-in devices." />
      {feedback ? <div className="notice">{feedback}</div> : null}
      <div className="settings-grid">
        <Card>
          <div className="card-title"><span className="card-title__icon"><KeyRound size={20} /></span><div><h3>Change password</h3><p className="muted">A recent reauthentication is required.</p></div></div>
          <form className="auth-form" onSubmit={(event) => { event.preventDefault(); reauth.mutate(); }}>
            <Input id="oldPassword" label="Verify your current password" type="password" value={oldPassword} onChange={(event) => setOldPassword(event.target.value)} required />
            {reauth.error ? <p className="form-error">{reauth.error.message}</p> : null}
            <Button type="submit" variant="secondary" loading={reauth.isPending}>Verify identity</Button>
          </form>
          <div className="sub-divider" />
          <form className="auth-form" onSubmit={submitPassword}>
            <Input id="changePassword" label="New password" type="password" minLength={12} value={newPassword} onChange={(event) => setNewPassword(event.target.value)} required />
            {change.error ? <p className="form-error">{change.error.message}</p> : null}
            <Button type="submit" loading={change.isPending}><Lock size={16} /> Change password</Button>
          </form>
        </Card>
        <Card>
          <div className="card-title"><span className="card-title__icon"><MonitorSmartphone size={20} /></span><div><h3>Active sessions</h3><p className="muted">Sign out devices you no longer use.</p></div></div>
          {sessions.isLoading ? <p className="muted">Loading sessions…</p> : sessions.error ? <p className="form-error">{sessions.error.message}</p> : (
            <div className="session-list">
              {sessions.data?.items.map((session) => (
                <div className="session-item" key={session.sessionId}>
                  <MonitorSmartphone size={20} />
                  <div><strong>{session.current ? 'This device' : 'Other device'}</strong><small>Created {new Date(session.createdAt).toLocaleDateString()}</small></div>
                  <Button variant="ghost" disabled={revoke.isPending} onClick={() => revoke.mutate(session.sessionId)}><Trash2 size={16} /></Button>
                </div>
              ))}
              {!sessions.data?.items.length ? <p className="muted">No sessions found.</p> : null}
              {revoke.error ? <p className="form-error">{revoke.error.message}</p> : null}
            </div>
          )}
          <div className="sub-divider" />
          <p className="muted security-footnote"><ShieldCheck size={16} /> Session tokens are kept in secure, HttpOnly browser cookies.</p>
        </Card>
      </div>
    </div>
  );
};
