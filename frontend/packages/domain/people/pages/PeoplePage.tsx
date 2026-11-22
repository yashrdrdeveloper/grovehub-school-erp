import { useState, type FormEvent } from 'react';
import { Plus, Search, UserRound, UserRoundX } from '@Components/icons';

import { personApi, type Person } from '@Business/foundation';
import { Button, Card, Input, PageHeading, StatusPill } from '@Components';
import { useMutation, useQuery, useQueryClient } from '@Components/hooks/query';

export const PeoplePage = () => {
  const queryClient = useQueryClient();
  const [searchText, setSearchText] = useState('');
  const [adding, setAdding] = useState(false);
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [key, setKey] = useState(() => crypto.randomUUID());

  const persons = useQuery({
    queryKey: ['persons', searchText],
    queryFn: () => personApi.search({
      searchText,
      pageInfo: { rows: 50 },
    }),
  });

  const create = useMutation({
    mutationFn: () => personApi.create(
      { firstName, lastName },
      { emails: [{ value: email, isPrimary: true }] },
      key,
    ),
    onSuccess: () => {
      setAdding(false);
      setMessage('Person created in Foundation.');
      setKey(crypto.randomUUID());
      setFirstName('');
      setLastName('');
      setEmail('');
      void queryClient.invalidateQueries({ queryKey: ['persons'] });
    },
  });

  const deactivate = useMutation({
    mutationFn: (gid: string) => personApi.deactivate(gid),
    onSuccess: () => void queryClient.invalidateQueries({ queryKey: ['persons'] }),
  });

  const submit = (event: FormEvent) => {
    event.preventDefault();
    create.mutate();
  };

  return (
    <div className="page">
      <PageHeading
        eyebrow="Foundation"
        title="People"
        description="People, their permanent GIDs, and contact information."
        action={<Button onClick={() => setAdding((value) => !value)}><Plus size={17} /> New person</Button>}
      />
      {message ? <div className="notice">{message}</div> : null}
      {adding ? (
        <Card>
          <h3>Create person</h3>
          <form className="form-grid" onSubmit={submit}>
            <Input id="personFirst" label="First name" value={firstName} onChange={(event) => setFirstName(event.target.value)} required />
            <Input id="personLast" label="Last name" value={lastName} onChange={(event) => setLastName(event.target.value)} required />
            <Input id="personEmail" label="Primary email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required />
            {create.error ? <p className="form-error">{create.error.message}</p> : null}
            <div className="form-actions"><Button type="button" variant="secondary" onClick={() => setAdding(false)}>Cancel</Button><Button type="submit" loading={create.isPending}>Create person</Button></div>
          </form>
        </Card>
      ) : null}
      <Card>
        <div className="toolbar">
          <div className="search-control"><Search size={17} /><input aria-label="Find a person" placeholder="Search people" value={searchText} onChange={(event) => setSearchText(event.target.value)} /></div>
          <Button variant="secondary" onClick={() => void persons.refetch()}>Refresh</Button>
        </div>
        {persons.isLoading ? <p className="table-feedback">Loading people…</p> : persons.error ? (
          <p className="form-error">{persons.error.message}</p>
        ) : !persons.data?.items.length ? (
          <p className="table-feedback">No people found.</p>
        ) : (
          <div className="table-scroll"><table className="data-table">
            <thead><tr><th>Person</th><th>GID</th><th>Status</th><th>Action</th></tr></thead>
            <tbody>
              {persons.data.items.map((person: Person) => (
                <tr key={person.gid}>
                  <td><div className="member-cell"><span className="member-avatar"><UserRound size={17} /></span><strong>{person.firstName} {person.lastName}</strong></div></td>
                  <td className="mono small">{person.gid}</td>
                  <td><StatusPill value={person.status} /></td>
                  <td>
                    <Button variant="ghost" disabled={deactivate.isPending || person.status !== 'ACTIVE'} onClick={() => {
                      if (window.confirm('Deactivate this person?')) deactivate.mutate(person.gid);
                    }}><UserRoundX size={16} /> Deactivate</Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table></div>
        )}
        {deactivate.error ? <p className="form-error">{deactivate.error.message}</p> : null}
      </Card>
    </div>
  );
};
