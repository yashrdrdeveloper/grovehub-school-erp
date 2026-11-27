import { useEffect, useState } from 'react';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { Input } from '@Components/components/Input';
import { __ } from '@Base/i18n';
import { ApiError, errorText, request, search } from '../../../business/api/client';
import { formatDate, parseObject } from '../../../business/api/helpers';
import type { Approval, Environment, Execution, Operator, Reconciliation, Template } from '../../../business/api/types';
import { ActionForm } from '../../../shell/components/ActionForm';
import { PagedRecords } from '../../../shell/components/PagedRecords';

function RequestForm({ close, success }: { close:()=>void; success:()=>void }) {
 const [templates,setTemplates]=useState<Template[]>([]);
 const [environments,setEnvironments]=useState<Environment[]>([]);
 const [templateId,setTemplateId]=useState('');
 const [error,setError]=useState('');
 const [ready,setReady]=useState(false);
 useEffect(()=>{
  const abort=new AbortController();
  void Promise.all([search<Template>('api-templates',1,100,abort.signal),search<Environment>('environments',1,100,abort.signal)]).then(([t,e])=>{if(!abort.signal.aborted){setTemplates(t.records);setEnvironments(e.records);setReady(true);}}).catch((failure)=>{if(!abort.signal.aborted)setError(errorText(failure));});
  return ()=>abort.abort();
 },[]);
 const template=templates.find((t)=>t.id===templateId);
 return <ActionForm title={__('Submit approval request')} onClose={close} onSuccess={success} submit={(data)=>{if(!ready||!template)throw new ApiError(400,'INVALID_REQUEST');return request('/internal/protected/requests','POST',{templateId,environmentId:data.get('environmentId'),schoolId:data.get('schoolId')||null,input:parseObject(String(data.get('input'))),reason:data.get('reason')});}}>
  {error && <p role="alert" className="error">{error}</p>}
  <label>{__('Template revision')}<select value={templateId} onChange={(event)=>setTemplateId(event.target.value)} required data-testid="request-template"><option value="">—</option>{templates.map((t)=><option key={t.id} value={t.id}>{t.code} · v{t.version}</option>)}</select></label>
  <label>{__('Environment')}<select name="environmentId" required data-testid="request-environment"><option value="">—</option>{environments.map((e)=><option key={e.id} value={e.id}>{e.name}</option>)}</select></label>
  {template?.scope==='SCHOOL' && <Input name="schoolId" label={__('School ID')} required dataTestId="request-school"/>}
  {template && <details><summary>{__('Input schema')}</summary><pre>{JSON.stringify(template.bodySchema,null,2)}</pre></details>}
  <label>{__('Input (JSON object)')}<textarea name="input" required defaultValue="{}" rows={5} data-testid="request-input"/></label>
  <label>{__('Reason')}<textarea name="reason" required maxLength={1000} rows={3} data-testid="request-reason"/></label>
 </ActionForm>;
}
function RequestDetail({ approval, operator, close, success }: { approval:Approval; operator:Operator; close:()=>void; success:()=>void }) {
 const [execution,setExecution]=useState<Execution|null>(null);
 const [review,setReview]=useState<Reconciliation|null>(null);
 const [error,setError]=useState('');
 useEffect(()=>{
  const abort=new AbortController();
  if(approval.status==='APPROVED'){
   void request<Execution>('/internal/protected/requests/'+approval.id+'/execution','GET',undefined,abort.signal).then((v)=>{if(!abort.signal.aborted)setExecution(v);}).catch((failure)=>{if(!abort.signal.aborted&&!(failure instanceof ApiError&&failure.status===404))setError(errorText(failure));});
   void request<Reconciliation>('/internal/protected/requests/'+approval.id+'/reconciliation','GET',undefined,abort.signal).then((v)=>{if(!abort.signal.aborted)setReview(v);}).catch((failure)=>{if(!abort.signal.aborted&&!(failure instanceof ApiError&&failure.status===404))setError(errorText(failure));});
  }
  return ()=>abort.abort();
 },[approval.id,approval.status]);
 const canReview=operator.persona==='L2'&&approval.status==='PENDING';
 const canCancel=operator.id===approval.requesterId&&approval.status==='PENDING';
 const canExecute=approval.status==='APPROVED'&&!execution&&(operator.persona==='L2'||operator.id===approval.requesterId);
 const canReconcile=operator.persona==='L2'&&!review&&(execution?.status==='UNKNOWN'||execution?.status==='RUNNING');
 const hasAction=canReview||canCancel||canExecute||canReconcile;
 const [action,setAction]=useState('');
 async function submit(data:FormData){
  const path='/internal/protected/requests/'+approval.id+'/';
  if(action==='execute')return request(path+action,'POST',{version:approval.version});
  if(action==='reconcile')return request(path+action,'POST',{outcome:data.get('outcome'),evidence:data.get('comment')});
  if(['approve','reject','cancel'].includes(action))return request(path+action,'POST',{version:approval.version,comment:data.get('comment')});
  throw new ApiError(400,'INVALID_REQUEST');
 }
 return <section className="action-card"><div className="page-heading"><h2>{__('Request detail')}</h2><Button variant={BUTTON_VARIANT.GHOST} onClick={close} dataTestId="request-detail-close">{__('Close')}</Button></div><span className="status">{approval.status}</span><p>{approval.reason}</p><details><summary>{__('Approved request data')}</summary><pre>{JSON.stringify(approval,null,2)}</pre></details>
 {error&&<p className="error" role="alert">{error}</p>}
 {execution&&<><h3>{__('Execution')}</h3><pre>{JSON.stringify(execution,null,2)}</pre></>}
 {review&&<><h3>{__('Reconciliation')}</h3><pre>{JSON.stringify(review,null,2)}</pre></>}
 {hasAction&&<ActionForm title={__('Request action')} onClose={close} onSuccess={success} submit={submit}><label>{__('Action')}<select value={action} onChange={(event)=>setAction(event.target.value)} required data-testid="request-action"><option value="">—</option>{canReview&&<><option value="approve">{__('Approve')}</option><option value="reject">{__('Reject')}</option></>}{canCancel&&<option value="cancel">{__('Cancel')}</option>}{canExecute&&<option value="execute">{__('Execute approved request')}</option>}{canReconcile&&<option value="reconcile">{__('Reconcile')}</option>}</select></label>
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
 <PagedRecords<Approval> resource={(inbox&&'approvals')||'requests'} title={title()} revision={revision} rowKey={(r)=>r.id} onSelect={setSelected} columns={[{title:__('Select'),cell:(r)=>inbox&&<input type="checkbox" checked={chosen.some((item)=>item.id===r.id)} onChange={()=>toggle(r)} aria-label={__('Select request')+' '+r.id} data-testid={'select-request-'+r.id}/>},{title:__('Reason'),cell:(r)=>r.reason},{title:__('Status'),cell:(r)=><span className="status">{r.status}</span>},{title:__('Version'),cell:(r)=>r.version},{title:__('Created'),cell:(r)=>formatDate(r.createdAt)}]}/>
 {selected&&<RequestDetail key={selected.id} approval={selected} operator={operator} close={close} success={success}/>}
 </>;
}
