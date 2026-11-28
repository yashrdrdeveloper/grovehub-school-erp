import { useState } from 'react';
import { NavLink } from '@Components/router';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import type { Operator } from '../../../business/api/types';
import { RequestForm } from '../../requests/pages/requestForm/RequestForm';
import { useRequestRegistry } from '../../requests/pages/requestForm/hooks/useRequestRegistry';
import type { RequestTarget } from '../../requests/pages/requestForm/types/requestForm.props';
import { buildSchoolRecords, filterSchools } from './helpers/schools.records';
import { SchoolDetails } from './components/SchoolDetails';
import type { SchoolAction, SchoolRecord } from './types/schools.types';

export function SchoolsPage({ operator }: { operator: Operator }) {
 const [revision, setRevision] = useState(0);
 const [query, setQuery] = useState('');
 const [environmentId, setEnvironmentId] = useState('');
 const [status, setStatus] = useState('');
 const [selectedKey, setSelectedKey] = useState('');
 const [target, setTarget] = useState<RequestTarget | null>(null);
 const [submitted, setSubmitted] = useState(false);
 const registry = useRequestRegistry(revision);
 const records = buildSchoolRecords(registry.grants, registry.environments, registry.templates);
 const visible = filterSchools(records, query, environmentId, status);
 const selected = records.find((record) => record.key === selectedKey);
 function refresh() {
  setTarget(null);
  setSelectedKey('');
  setRevision((current) => current + 1);
 }
 function requestChange(record: SchoolRecord, action: SchoolAction) {
  if (!record.snapshotReady || !record.status || !record.actions.includes(action)) return;
  setSubmitted(false);
  setTarget({ templateId: action.template.id, environmentId: record.environment.id, schoolId: record.schoolId, schoolVersion: record.version, schoolStatus: record.status, schoolName: record.name });
 }
 function success() {
  setSubmitted(true);
  refresh();
 }
 return <>
  <header className="page-heading" data-testid="schools-heading"><div><span className="eyebrow">{__('REGISTRY')}</span><h1>{__('Schools')}</h1><p className="muted">{__('View schools within your configured grants and request approved lifecycle changes.')}</p></div><Button variant={BUTTON_VARIANT.OUTLINE} onClick={refresh} disabled={Boolean(target) || (!registry.ready && !registry.error)} dataTestId="schools-refresh">{__('Refresh')}</Button></header>
  <nav className="inline" aria-label={__('School workflow links')}><NavLink className="school-link" to="/requests">{__('My requests')}</NavLink>{operator.persona === 'L2' && <><NavLink className="school-link" to="/approvals">{__('Approval inbox')}</NavLink><NavLink className="school-link" to="/audit/schools">{__('School mutation history (all schools)')}</NavLink></>}</nav>
  {submitted && <p role="status" className="action-card">{__('Approval request submitted. School status has not changed yet.')}</p>}
  {registry.error && <p role="alert" className="error">{registry.error}</p>}
  {!registry.ready && !registry.error && <p role="status">{__('Loading authorized schools…')}</p>}
  {registry.ready && <>
   <section className="action-card" data-testid="schools-filters"><div className="school-filters">
    <label>{__('Search schools')}<input value={query} onChange={(event) => setQuery(event.target.value)} placeholder={__('Name, code or ID')} data-testid="schools-search"/></label>
    <label>{__('Environment')}<select value={environmentId} onChange={(event) => setEnvironmentId(event.target.value)} data-testid="schools-environment"><option value="">{__('All authorized environments')}</option>{registry.environments.filter((item) => records.some((record) => record.environment.id === item.id)).map((item) => <option value={item.id} key={item.id}>{item.name}</option>)}</select></label>
    <label>{__('Status')}<select value={status} onChange={(event) => setStatus(event.target.value)} data-testid="schools-status"><option value="">{__('All statuses')}</option><option value="ACTIVE">{__('Active')}</option><option value="DISABLED">{__('Disabled')}</option></select></label>
   </div><p className="muted small" role="status">{visible.length} {__('matching school and environment records')}</p></section>
   {!target && selected && <SchoolDetails record={selected} canViewHistory={operator.persona === 'L2'} close={() => setSelectedKey('')} requestChange={requestChange}/>}
   {target && <RequestForm key={[target.templateId,target.environmentId,target.schoolId,target.schoolVersion].join(':')} target={target} close={() => setTarget(null)} success={success}/>}
   {!visible.length && <p role="status" className="action-card">{__('No matching schools in your authorized grants.')}</p>}
   <section className="school-list" aria-label={__('Authorized schools')} data-testid="schools-list">{visible.map((record) => <article className="school-record" key={record.key}>
    <div><h3>{record.name}</h3><p className="muted small">{record.code || record.schoolId}</p></div>
    <div><span className="status">{record.snapshotReady && record.status || __('Unavailable')}</span><p className="muted small">{record.environment.name}</p></div>
    <Button variant={BUTTON_VARIANT.OUTLINE} disabled={Boolean(target)} onClick={() => setSelectedKey(record.key)} dataTestId={'school-open-'+record.environment.id+'-'+record.schoolId}>{__('Details')}</Button>
   </article>)}</section>
  </>}
 </>;
}
