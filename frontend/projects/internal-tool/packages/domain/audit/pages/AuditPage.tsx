import type { WorkflowEvent } from '../../../business/api/types';
import { formatDate } from '../../../business/api/helpers';
import { PagedRecords } from '../../../shell/components/pagedRecords/PagedRecords';

export function AuditPage() {
 return <><header className="page-heading"><div><span className="eyebrow">{__('ACCOUNTABILITY')}</span><h1>{__('Workflow history')}</h1><p className="muted">{__('Internal Tool activity recorded in its own database.')}</p></div></header><PagedRecords<WorkflowEvent> resource="workflow-events" title={__('Workflow events')} revision={0} rowKey={(r)=>r.id} columns={[{id:'action',title:__('Action'),cell:(r)=>r.action},{id:'actor',title:__('Actor'),cell:(r)=>r.actorId},{id:'resource',title:__('Resource'),cell:(r)=>r.resourceId},{id:'occurred-at',title:__('Occurred at'),cell:(r)=>formatDate(r.occurredAt)}]}/></>;
}
