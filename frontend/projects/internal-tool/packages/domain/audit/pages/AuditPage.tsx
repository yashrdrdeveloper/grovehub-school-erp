import { __ } from '@Base/i18n';
import type { WorkflowEvent } from '../../../business/api/types';
import { formatDate } from '../../../business/api/helpers';
import { PagedRecords } from '../../../shell/components/PagedRecords';

export function AuditPage() {
 return <><header className="page-heading"><div><span className="eyebrow">{__('ACCOUNTABILITY')}</span><h1>{__('Workflow history')}</h1><p className="muted">{__('Internal Tool activity recorded in its own database.')}</p></div></header><PagedRecords<WorkflowEvent> resource="workflow-events" title={__('Workflow events')} revision={0} rowKey={(r)=>r.id} columns={[{title:__('Action'),cell:(r)=>r.action},{title:__('Actor'),cell:(r)=>r.actorId},{title:__('Resource'),cell:(r)=>r.resourceId},{title:__('Occurred at'),cell:(r)=>formatDate(r.occurredAt)}]}/></>;
}
