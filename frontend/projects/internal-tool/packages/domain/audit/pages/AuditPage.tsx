import { useState } from 'react';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import type { WorkflowEvent } from '../../../business/api/types';
import { request } from '../../../business/api/client';
import { formatDate } from '../../../business/api/helpers';
import { PagedRecords } from '../../../shell/components/pagedRecords/PagedRecords';
import { ActionForm } from '../../../shell/components/ActionForm';
import { useAuditDelivery } from './hooks/useAuditDelivery';
import type { BlockedEvent, OperatorActivity } from './types/audit.types';

type View = 'local' | 'central' | 'blocked';

export function AuditPage() {
  const [view, setView] = useState<View>('local');
  const [revision, setRevision] = useState(0);
  const [selected, setSelected] = useState<BlockedEvent | null>(null);
  const { summary, error } = useAuditDelivery(revision);
  function changeView(next: View) {
    setSelected(null);
    setView(next);
  }
  function refresh() { setRevision((current) => current + 1); }
  return <>
    <header className="page-heading" data-testid="audit-heading">
      <div><span className="eyebrow">{__('ACCOUNTABILITY')}</span><h1>{__('Audit history')}</h1><p className="muted">{__('Inspect local history, centralized operator activities and delivery status.')}</p></div>
      <Button variant={BUTTON_VARIANT.OUTLINE} onClick={refresh} dataTestId="audit-refresh">{__('Refresh')}</Button>
    </header>
    <section className="action-card" data-testid="audit-delivery-summary">
      <h2>{__('Audit delivery')}</h2>
      {error && <p className="error" role="alert">{error}</p>}
      {!summary && !error && <p role="status">{__('Loading delivery status…')}</p>}
      {summary && <>
        <span className="status">{summary.enabled && __('Delivery enabled')}{!summary.enabled && __('Delivery not configured')}</span>
        <div className="inline"><span>{__('Pending')}: {summary.pending}</span><span>{__('Blocked')}: {summary.blocked}</span><span>{__('Delivered')}: {summary.delivered}</span></div>
        <p className="muted">{__('Oldest pending')}: {formatDate(summary.oldestPendingAt)} · {__('Last delivered')}: {formatDate(summary.lastDeliveredAt)}</p>
      </>}
    </section>
    <nav className="inline" aria-label={__('Audit views')} data-testid="audit-views">
      <Button variant={BUTTON_VARIANT.OUTLINE} aria-pressed={view==='local'} onClick={()=>changeView('local')} dataTestId="audit-local">{__('Local history')}</Button>
      <Button variant={BUTTON_VARIANT.OUTLINE} aria-pressed={view==='central'} onClick={()=>changeView('central')} dataTestId="audit-central">{__('Centralized history')}</Button>
      <Button variant={BUTTON_VARIANT.OUTLINE} aria-pressed={view==='blocked'} onClick={()=>changeView('blocked')} dataTestId="audit-blocked">{__('Blocked delivery')}</Button>
    </nav>
    {view==='local' && <PagedRecords<WorkflowEvent> key="local" resource="workflow-events" title={__('Local workflow events')} revision={revision} rowKey={(r)=>r.id} columns={[{id:'action',title:__('Action'),cell:(r)=>r.action},{id:'actor',title:__('Actor'),cell:(r)=>r.actorId},{id:'resource',title:__('Resource'),cell:(r)=>r.resourceId},{id:'correlation',title:__('Correlation ID'),cell:(r)=>r.correlationId||'—'},{id:'occurred-at',title:__('Occurred at'),cell:(r)=>formatDate(r.occurredAt)}]}/>}
    {view==='central' && <PagedRecords<OperatorActivity> key="central" resource="operator-activities" title={__('Centralized operator activities')} revision={revision} rowKey={(r)=>r.eventId} columns={[{id:'operation',title:__('Operation'),cell:(r)=>r.operation},{id:'actor',title:__('Actor'),cell:(r)=>r.actorId},{id:'resource',title:__('Resource'),cell:(r)=>r.resourceId},{id:'correlation',title:__('Correlation ID'),cell:(r)=>r.correlationId||'—'},{id:'occurred-at',title:__('Occurred at'),cell:(r)=>formatDate(r.occurredAt)}]}/>}
    {view==='blocked' && <PagedRecords<BlockedEvent> key="blocked" resource="audit-blocked" title={__('Blocked events')} revision={revision} rowKey={(r)=>r.eventId} onSelect={setSelected} columns={[{id:'action',title:__('Action'),cell:(r)=>r.action},{id:'actor',title:__('Actor'),cell:(r)=>r.actorId},{id:'resource',title:__('Resource'),cell:(r)=>r.resourceId},{id:'attempts',title:__('Attempts'),cell:(r)=>r.attempts},{id:'blocked-at',title:__('Blocked at'),cell:(r)=>formatDate(r.blockedAt)}]}/>}
    {selected && view==='blocked' && <ActionForm title={__('Requeue blocked occurrence')} onClose={()=>setSelected(null)} onSuccess={()=>{setSelected(null);refresh();}} submit={()=>request('/internal/protected/audit-blocked/'+selected.eventId+'/retry','POST')}>
      <p>{__('Investigate the delivery failure before retrying. The original event ID and content will not change.')}</p>
      <p>{selected.action} · {selected.eventId}</p>
    </ActionForm>}
  </>;
}
