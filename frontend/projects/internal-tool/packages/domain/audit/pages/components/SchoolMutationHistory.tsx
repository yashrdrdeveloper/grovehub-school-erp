import { useState, type FormEvent } from 'react';
import { useSearchParams } from '@Components/router';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { PagedRecords } from '../../../../shell/components/pagedRecords/PagedRecords';
import { formatDate } from '../../../../business/api/helpers';
import type { SchoolReceipt } from '../types/audit.types';
import { buildReceiptFilters, SCHOOL_MUTATION_OPERATIONS } from '../helpers/schoolMutation.filters';

function MutationFilter({ schoolId, operation, apply }: { schoolId: string; operation: string; apply: (id: string, op: string) => void }) {
 const [id, setId] = useState(schoolId);
 const [kind, setKind] = useState(operation);
 const valid = buildReceiptFilters(id, kind) !== null;
 function submit(event: FormEvent<HTMLFormElement>) {
  event.preventDefault();
  if (valid) apply(id,kind);
 }
 return <section className="action-card" data-testid="school-history-filters">
  <form className="school-filters" onSubmit={submit}>
   <label>{__('School ID')}<input value={id} onChange={(event) => setId(event.target.value)} maxLength={36} placeholder={__('All schools')} data-testid="history-school-id"/></label>
   <label>{__('Operation')}<select value={kind} onChange={(event) => setKind(event.target.value)} data-testid="history-operation"><option value="">{__('All operations')}</option>{SCHOOL_MUTATION_OPERATIONS.map((item) => <option key={item} value={item}>{item}</option>)}</select></label>
   <Button type="submit" disabled={!valid} dataTestId="history-apply">{__('Apply filters')}</Button>
   <Button type="button" variant={BUTTON_VARIANT.OUTLINE} onClick={() => { setId(''); setKind(''); apply('',''); }} dataTestId="history-reset">{__('Reset')}</Button>
  </form>
  {!valid && <p role="alert" className="error">{__('Enter a canonical school UUID and a valid operation.')}</p>}
  <p className="muted small">{__('Only committed receipts delivered to central Audit appear here. History is scoped by school ID, not environment.')}</p>
 </section>;
}

export function SchoolMutationHistory({ revision }: { revision: number }) {
 const [params, setParams] = useSearchParams();
 const schoolId = params.get('schoolId') ?? '';
 const operation = params.get('operation') ?? '';
 const filters = buildReceiptFilters(schoolId, operation);
 const key = JSON.stringify([schoolId, operation]);
 function apply(id: string, op: string) {
  const next = buildReceiptFilters(id,op);
  if (next) setParams(next);
 }
 return <>
  <MutationFilter key={key} schoolId={schoolId} operation={operation} apply={apply}/>
  {filters === null && <p role="alert" className="error">{__('Invalid history filters. Reset or correct them before searching.')}</p>}
  {filters !== null && <PagedRecords<SchoolReceipt> key={key} resource="school-mutation-receipts" filters={filters} title={__('School mutation receipts')} revision={revision} rowKey={(r)=>r.executionId} columns={[
   {id:'operation',title:__('Operation'),cell:(r)=>r.operation},
   {id:'execution',title:__('Execution ID'),cell:(r)=>r.executionId},
   {id:'school',title:__('School ID'),cell:(r)=>r.schoolId},
   {id:'version',title:__('Previous version'),cell:(r)=>r.previousVersion},
   {id:'correlation',title:__('Correlation ID'),cell:(r)=>r.correlationId||'—'},
   {id:'occurred-at',title:__('Occurred at'),cell:(r)=>formatDate(r.occurredAt)}
  ]}/>}
 </>;
}
