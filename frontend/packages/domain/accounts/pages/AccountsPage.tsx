import { useMemo, useState, type FormEvent } from 'react';
import { Plus, RefreshCw, Search, UserRound } from '@Components/icons';
import { useMutation, useQuery, useQueryClient } from '@Components/hooks/query';
import type { ColumnDef } from '@tanstack/react-table';
import { Link } from 'react-router-dom';

import { accountsApi, type CreateAccount } from '@Business/accounts';
import type { Account } from '@Business/auth';
import { Button, Card, DataTable, Input, PageHeading, StatusPill } from '@Components';

const generateKey = (): string => crypto.randomUUID();

export const AccountsPage = () => {
  const queryClient = useQueryClient();
  const [text, setText] = useState('');
  const [creating, setCreating] = useState(false);
  const [draft, setDraft] = useState<CreateAccount>({ person: { firstName: '', lastName: '' }, email: '' });
  const [idempotencyKey, setIdempotencyKey] = useState(generateKey);
  const [feedback, setFeedback] = useState('');

  const accounts = useQuery({
    queryKey: ['accounts', text],
    queryFn: () => accountsApi.search({
      ...(text.trim() ? { searchText: text.trim() } : {}),
      pageInfo: { rows: 50 },
    }),
  });

  const create = useMutation({
    mutationFn: () => accountsApi.create(draft, idempotencyKey),
    onSuccess: () => {
      setCreating(false);
      setFeedback('Account created. Invitation queued for delivery.');
      setDraft({ person: { firstName: '', lastName: '' }, email: '' });
      setIdempotencyKey(generateKey());
      void queryClient.invalidateQueries({ queryKey: ['accounts'] });
    },
  });

  const columns = useMemo<ColumnDef<Account, unknown>[]>(() => [
    {
      id: 'member',
      header: 'Account',
      cell: ({ row }) => (
        <div className="member-cell">
          <span className="member-avatar"><UserRound size={16} /></span>
          <span><strong><Link to={`/accounts/${row.original.accountId}`}>{row.original.username ?? row.original.email.split('@')[0]}</Link></strong><small>{row.original.email}</small></span>
        </div>
      ),
    },
    { accessorKey: 'personGid', header: 'Person GID', cell: ({ getValue }) => <span className="mono small">{String(getValue()).slice(0, 16)}…</span> },
    { accessorKey: 'status', header: 'Status', cell: ({ getValue }) => <StatusPill value={String(getValue())} /> },
    { accessorKey: 'createdAt', header: 'Created', cell: ({ getValue }) => new Date(String(getValue())).toLocaleDateString() },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }) => (
        <Button variant="ghost" onClick={() => {
          void accountsApi.resendInvitation(row.original.accountId)
            .then(() => setFeedback('Invitation queued.'))
            .catch((error: Error) => setFeedback(error.message));
        }}>Resend invite</Button>
      ),
    },
  ], []);

  const submit = (event: FormEvent) => {
    event.preventDefault();
    setFeedback('');
    create.mutate();
  };

  return (
    <div className="page">
      <PageHeading
        eyebrow="Identity management"
        title="Accounts"
        description="School staff access, account setup invitations and identity lifecycle."
        action={<Button onClick={() => setCreating((value) => !value)}><Plus size={17} /> New account</Button>}
      />
      {feedback ? <div className="notice">{feedback}</div> : null}
      {creating ? (
        <Card className="form-card">
          <h3>Invite a school member</h3>
          <p className="muted">We'll create or link their Foundation Person GID and queue a secure setup invitation.</p>
          <form className="form-grid" onSubmit={submit}>
            <Input id="firstName" label="First name" value={draft.person.firstName} required onChange={(event) => setDraft({ ...draft, person: { ...draft.person, firstName: event.target.value } })} />
            <Input id="lastName" label="Last name" value={draft.person.lastName} required onChange={(event) => setDraft({ ...draft, person: { ...draft.person, lastName: event.target.value } })} />
            <Input id="accountEmail" label="Email" type="email" required value={draft.email} onChange={(event) => setDraft({ ...draft, email: event.target.value })} />
            <Input id="phoneNumber" label="Phone (optional)" type="tel" value={draft.phoneNumber ?? ''} onChange={(event) => setDraft({ ...draft, phoneNumber: event.target.value })} />
            {create.error ? <p className="form-error">{create.error.message}</p> : null}
            <div className="form-actions"><Button type="button" variant="secondary" onClick={() => setCreating(false)}>Cancel</Button><Button type="submit" loading={create.isPending}>Create & invite</Button></div>
          </form>
        </Card>
      ) : null}
      <Card>
        <div className="toolbar">
          <div className="search-control"><Search size={17} /><input aria-label="Search accounts" placeholder="Search names, email or username" value={text} onChange={(event) => setText(event.target.value)} /></div>
          <Button variant="secondary" onClick={() => void accounts.refetch()}><RefreshCw size={16} /> Refresh</Button>
        </div>
        {accounts.isLoading ? <div className="table-feedback">Loading accounts…</div> : accounts.error ? (
          <div className="table-feedback form-error" role="alert">{accounts.error.message}</div>
        ) : <DataTable rows={accounts.data?.items ?? []} columns={columns} emptyTitle="No accounts found" />}
        {accounts.data?.pageInfo?.hasNext ? <p className="table-footer">More results are available. Cursor pagination is supported by the API client.</p> : null}
      </Card>
    </div>
  );
};
