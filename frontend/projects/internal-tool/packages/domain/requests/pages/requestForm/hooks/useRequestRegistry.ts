import { useEffect, useState } from 'react';
import { errorText, search } from '../../../../../business/api/client';
import type { Environment, Template } from '../../../../../business/api/types';

async function loadRegistry<T>(resource: string, signal: AbortSignal, pageNumber = 1, records: T[] = []): Promise<T[]> {
  const page = await search<T>(resource, pageNumber, 100, signal);
  const next = [...records, ...page.records];
  if (!page.hasMore) return next;
  // Bound selector memory; use registry search when the catalog grows beyond this.
  if (pageNumber >= 10) throw new Error('Registry exceeds selector limit');
  return loadRegistry(resource, signal, pageNumber + 1, next);
}

export function useRequestRegistry() {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [error, setError] = useState('');
  const [ready, setReady] = useState(false);
  useEffect(() => {
    const abort = new AbortController();
    void Promise.all([
      loadRegistry<Template>('api-templates', abort.signal),
      loadRegistry<Environment>('environments', abort.signal),
    ]).then(([templateRows, environmentRows]) => {
      if (abort.signal.aborted) return;
      setTemplates(templateRows);
      setEnvironments(environmentRows);
      setReady(true);
    }).catch((failure: unknown) => {
      if (!abort.signal.aborted) setError(errorText(failure));
    });
    return () => abort.abort();
  }, []);
  return { templates, environments, error, ready };
}
