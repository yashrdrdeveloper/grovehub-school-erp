import { useEffect, useState } from 'react';
import { errorText, request } from '../../../../business/api/client';
import type { DeliverySummary } from '../types/audit.types';

export function useAuditDelivery(revision: number) {
  const [summary, setSummary] = useState<DeliverySummary | null>(null);
  const [error, setError] = useState('');
  useEffect(() => {
    const controller = new AbortController();
    setError('');
    void request<DeliverySummary>('/internal/protected/audit-delivery/summary', 'GET', undefined, controller.signal)
      .then((value) => { if (!controller.signal.aborted) setSummary(value); })
      .catch((failure: unknown) => {
        if (!controller.signal.aborted) { setSummary(null); setError(errorText(failure)); }
      });
    return () => controller.abort();
  }, [revision]);
  return { summary, error };
}
