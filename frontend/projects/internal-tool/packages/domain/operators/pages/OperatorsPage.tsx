import { useState } from 'react';
import { Button } from '@Components/components/Button';
import { Input } from '@Components/components/Input';
import { __ } from '@Base/i18n';
import type { Operator } from '../../../business/api/types';
import { request } from '../../../business/api/client';
import { ActionForm } from '../../../shell/components/ActionForm';
import { PagedRecords } from '../../../shell/components/PagedRecords';

export function OperatorsPage({ operator }: { operator: Operator }) {
  const [create,setCreate]=useState(false);
  const [selected,setSelected]=useState<Operator|null>(null);
  const [revision,setRevision]=useState(0);
  function close(){setCreate(false);setSelected(null);}
  function success(){close();setRevision((n)=>n+1);}
  if(operator.persona!=='L2') return <p role="alert">{__('Access denied.')}</p>;
  return <><header className="page-heading"><div><span className="eyebrow">{__('ACCESS MANAGEMENT')}</span><h1>{__('Operators')}</h1><p className="muted">{__('Manage Internal Tool operator accounts and personas.')}</p></div><Button onClick={()=>{setSelected(null);setCreate(true);}} dataTestId="operator-create">{__('Create operator')}</Button></header>
    {create && <ActionForm title={__('Create operator')} onClose={close} onSuccess={success} submit={(data)=>request('/internal/protected/operators','POST',{email:data.get('email'),password:data.get('password'),persona:data.get('persona')})}><Input name="email" label={__('Email')} type="email" required dataTestId="operator-email"/><Input name="password" label={__('Initial password')} type="password" minLength={12} maxLength={256} autoComplete="new-password" required dataTestId="operator-password"/><label>{__('Persona')}<select name="persona" data-testid="operator-persona"><option>L1</option><option>L2</option></select></label></ActionForm>}
    <PagedRecords<Operator> resource="operators" title={__('Operators')} revision={revision} rowKey={(r)=>r.id} onSelect={(r)=>{setCreate(false);setSelected(r);}} columns={[{title:__('Email'),cell:(r)=>r.email},{title:__('Persona'),cell:(r)=>r.persona},{title:__('Status'),cell:(r)=><span className="status">{r.status}</span>}]}/>
    {selected && <ActionForm title={selected.email} onClose={close} onSuccess={success} submit={(data)=>request('/internal/protected/operators/'+selected.id+'/'+String(data.get('action')),'POST')}>
      {selected.id === operator.id && <p>{__('You cannot change your own status.')}</p>}
      <label>{__('Action')}<select name="action" required data-testid="operator-action"><option value="">—</option>{selected.id!==operator.id && selected.status==='ACTIVE' && <option value="disable">{__('Disable')}</option>}{selected.id!==operator.id && selected.status==='DISABLED' && <option value="activate">{__('Activate')}</option>}</select></label>
    </ActionForm>}
  </>;
}
