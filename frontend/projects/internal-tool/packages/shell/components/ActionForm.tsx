import { useState, type FormEvent, type ReactNode } from 'react';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { errorText } from '../../business/api/client';

export function ActionForm({ title, children, submit, onClose, onSuccess, submitDisabled = false }: { title: string; children: ReactNode; submit: (data: FormData) => Promise<unknown>; onClose: () => void; onSuccess: () => void; submitDisabled?: boolean }) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busy || submitDisabled) return;
    setBusy(true); setError('');
    try { await submit(new FormData(event.currentTarget)); onSuccess(); }
    catch (failure) { setError(errorText(failure)); }
    finally { setBusy(false); }
  }
  return <section className="action-card" data-testid="internal-action-form"><div className="page-heading"><h2>{title}</h2><Button variant={BUTTON_VARIANT.GHOST} onClick={onClose} disabled={busy} dataTestId="internal-action-close">{__('Close')}</Button></div>
    <form className="form-stack" onSubmit={handleSubmit}>{children}{error && <p className="error" role="alert">{error}</p>}<Button type="submit" loading={busy} disabled={submitDisabled} dataTestId="internal-action-submit">{__('Submit')}</Button></form>
  </section>;
}
