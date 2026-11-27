import { useEffect, useState } from 'react';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { Checkbox } from '@Components/components/Checkbox';
import { RequestForm } from './requestForm/RequestForm';
import { ApiError, errorText, request } from '../../../business/api/client';
import { formatDate } from '../../../business/api/helpers';
import type { Approval, Execution, Operator, Reconciliation, RequestCapabilities } from '../../../business/api/types';
import { ActionForm } from '../../../shell/components/ActionForm';
import { PagedRecords } from '../../../shell/components/pagedRecords/PagedRecords';

function RequestDetail({ approval, close, success }: { approval:Approval; close:()=>void; success:()=>void }) {
 const [capabilities,setCapabilities]=useState<RequestCapabilities|null>(null);
 const [execution,setExecution]=useState<Execution|null>(null);
 const [review,setReview]=useState<Reconciliation|null>(null);
 const [error,setError]=useState('');
 useEffect(()=>{
  const abort=new AbortController();
  setCapabilities(null);
  setExecution(null);
  setReview(null);
  setError('');
  void request<RequestCapabilities>('/internal/protected/requests/'+approval.id+'/capabilities','GET',undefined,abort.signal)
    .then((value)=>{if(!abort.signal.aborted)setCapabilities(value);})
    .catch((failure:unknown)=>{if(!abort.signal.aborted)setError(errorText(failure));});
  if(approval.status==='APPROVED'){
   void request<Execution>('/internal/protected/requests/'+approval.id+'/execution','GET',undefined,abort.signal).then((v)=>{if(!abort.signal.aborted)setExecution(v);}).catch((failure)=>{if(!abort.signal.aborted&&!(failure instanceof ApiError&&failure.status===404))setError(errorText(failure));});
   void request<Reconciliation>('/internal/protected/requests/'+approval.id+'/reconciliation','GET',undefined,abort.signal).then((v)=>{if(!abort.signal.aborted)setReview(v);}).catch((failure)=>{if(!abort.signal.aborted&&!(failure instanceof ApiError&&failure.status===404))setError(errorText(failure));});
  }
  return ()=>abort.abort();
 },[approval.id,approval.status,approval.version]);
 const current=capabilities?.version===approval.version;
 const canApprove=Boolean(current&&capabilities?.canApprove);
 const canReject=Boolean(current&&capabilities?.canReject);
 const canCancel=Boolean(current&&capabilities?.canCancel);
 const canExecute=Boolean(current&&capabilities?.canExecute&&!execution);
 const canReconcile=Boolean(current&&capabilities?.canReconcile&&!review&&execution);
 const hasAction=canApprove||canReject||canCancel||canExecute||canReconcile;
 const [action,setAction]=useState('');
 async function submit(data:FormData){
  const path='/internal/protected/requests/'+approval.id+'/';
  if(action==='execute'&&canExecute)return request(path+action,'POST',{version:approval.version});
  if(action==='reconcile'&&canReconcile)return request(path+action,'POST',{outcome:data.get('outcome'),evidence:data.get('comment')});
  if((action==='approve'&&canApprove)||(action==='reject'&&canReject)||(action==='cancel'&&canCancel))return request(path+action,'POST',{version:approval.version,comment:data.get('comment')});
  throw new ApiError(400,'INVALID_REQUEST');
 }
 return <section className="action-card"><div className="page-heading"><h2>{__('Request detail')}</h2><Button variant={BUTTON_VARIANT.GHOST} onClick={close} dataTestId="request-detail-close">{__('Close')}</Button></div><span className="status">{approval.status}</span><p>{approval.reason}</p><details><summary>{__('Approved request data')}</summary><pre>{JSON.stringify(approval,null,2)}</pre></details>
 {!capabilities&&!error&&<p role="status">{__('Checking permitted actions…')}</p>}
 {capabilities&&!current&&<p role="status">{__('Request changed. Close and refresh the list.')}</p>}
 {error&&<p className="error" role="alert">{error}</p>}
 {execution&&<><h3>{__('Execution')}</h3><pre>{JSON.stringify(execution,null,2)}</pre></>}
 {review&&<><h3>{__('Reconciliation')}</h3><pre>{JSON.stringify(review,null,2)}</pre></>}
 {hasAction&&<ActionForm title={__('Request action')} onClose={close} onSuccess={success} submit={submit}><label>{__('Action')}<select value={action} onChange={(event)=>setAction(event.target.value)} required data-testid="request-action"><option value="">—</option>{canApprove&&<option value="approve">{__('Approve')}</option>}{canReject&&<option value="reject">{__('Reject')}</option>}{canCancel&&<option value="cancel">{__('Cancel')}</option>}{canExecute&&<option value="execute">{__('Execute approved request')}</option>}{canReconcile&&<option value="reconcile">{__('Reconcile')}</option>}</select></label>
 {action==='reconcile'&&<label>{__('Outcome')}<select name="outcome" data-testid="reconciliation-outcome"><option>CONFIRMED_SUCCEEDED</option><option>CONFIRMED_FAILED</option><option>UNRESOLVED</option></select></label>}
 {action!=='execute'&&<label>{__('Comment or evidence')}<textarea name="comment" required={action==='reject'||action==='reconcile'} maxLength={1000} rows={3} data-testid="request-comment"/></label>}
 </ActionForm>}</section>;
}
export function RequestsPage({ inbox, operator }: { inbox:boolean; operator:Operator }) {
 const [create,setCreate]=useState(false);
 const [selected,setSelected]=useState<Approval|null>(null);
 const [revision,setRevision]=useState(0);
 const [bulk,setBulk]=useState(false);
 const [chosen,setChosen]=useState<Approval[]>([]);
 function close(){setCreate(false);setSelected(null);setBulk(false);}
 function success(){close();setChosen([]);setRevision((n)=>n+1);}
 function title(){if(inbox)return __('Approval inbox');return __('My requests');}
 function toggle(r:Approval){setChosen((old)=>{if(old.some((item)=>item.id===r.id))return old.filter((item)=>item.id!==r.id);if(old.length>=50)return old;return [...old,r];});}
 if(inbox&&operator.persona!=='L2')return <p role="alert">{__('Access denied.')}</p>;
 return <><header className="page-heading"><div><span className="eyebrow">{__('OPERATIONS')}</span><h1>{title()}</h1><p className="muted">{__('Review requests, decisions and approved executions.')}</p></div><div className="inline">{!inbox&&<Button onClick={()=>setCreate(true)} dataTestId="request-create">{__('New request')}</Button>}{inbox&&<Button disabled={!chosen.length} onClick={()=>setBulk(true)} dataTestId="bulk-approve">{__('Approve selected')} ({chosen.length})</Button>}</div></header>
 {create&&!inbox&&<RequestForm close={close} success={success}/>}
 {bulk&&inbox&&<ActionForm title={__('Confirm bulk approval')} onClose={close} onSuccess={success} submit={(data)=>request('/internal/protected/approvals/bulk-approve','POST',{items:chosen.map((r)=>({id:r.id,version:r.version})),comment:data.get('comment')})}><p>{chosen.length} {__('pending requests selected')}</p><label>{__('Comment')}<textarea name="comment" maxLength={1000} rows={3} data-testid="bulk-comment"/></label></ActionForm>}
 <PagedRecords<Approval> resource={(inbox&&'approvals')||'requests'} title={title()} revision={revision} rowKey={(r)=>r.id} onSelect={setSelected} columns={[{id:'select',title:__('Select'),cell:(r)=>inbox&&<Checkbox checked={chosen.some((item)=>item.id===r.id)} onChange={()=>toggle(r)} label={__('Select request {{id}}', { id: r.id })} dataTestId={'select-request-'+r.id}/>},{id:'reason',title:__('Reason'),cell:(r)=>r.reason},{id:'status',title:__('Status'),cell:(r)=><span className="status">{r.status}</span>},{id:'version',title:__('Version'),cell:(r)=>r.version},{id:'created',title:__('Created'),cell:(r)=>formatDate(r.createdAt)}]}/>
 {selected&&<RequestDetail key={selected.id} approval={selected} close={close} success={success}/>}
 </>;
}
