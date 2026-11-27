import { useState } from 'react';
import { Button } from '@Components/components/Button';
import { Input } from '@Components/components/Input';
import { __ } from '@Base/i18n';
import { request } from '../../../business/api/client';
import { parseObject } from '../../../business/api/helpers';
import type { Environment, Operator, Template } from '../../../business/api/types';
import { ActionForm } from '../../../shell/components/ActionForm';
import { PagedRecords } from '../../../shell/components/PagedRecords';

export function RegistryPage({ kind, operator }: { kind: 'environments' | 'api-templates'; operator: Operator }) {
  const [create, setCreate] = useState(false);
  const [revision, setRevision] = useState(0);
  const [selected, setSelected] = useState<Template | null>(null);
  const environments = kind === 'environments';
  function close() { setCreate(false); }
  function success() { close(); setRevision((n) => n + 1); }
  function title() { if (environments) return __('Environments'); return __('API Explorer'); }
  return <><header className="page-heading"><div><span className="eyebrow">{__('REGISTRY')}</span><h1>{title()}</h1><p className="muted">{__('Manage logical environments and immutable restricted API definitions.')}</p></div>{operator.persona === 'L2' && <Button onClick={() => setCreate(true)} dataTestId="registry-create">{__('Create')}</Button>}</header>
    {create && environments && <ActionForm title={__('Create environment')} onClose={close} onSuccess={success} submit={(data) => request('/internal/protected/environments','POST',{code: data.get('code'), name: data.get('name')})}><Input name="code" label={__('Code')} required dataTestId="environment-code"/><Input name="name" label={__('Name')} required dataTestId="environment-name"/></ActionForm>}
    {create && !environments && <ActionForm title={__('Register template revision')} onClose={close} onSuccess={success} submit={(data) => request('/internal/protected/api-templates','POST',{code:data.get('code'),module:data.get('module'),method:data.get('method'),endpoint:data.get('endpoint'),scope:data.get('scope'),headers:parseObject(String(data.get('headers'))),bodySchema:parseObject(String(data.get('schema')))})}>
      <Input name="code" label={__('Template code')} required dataTestId="template-code"/><Input name="module" label={__('Module')} required dataTestId="template-module"/>
      <label>{__('Method')}<select name="method" data-testid="template-method">{['GET','POST','PATCH','DELETE'].map((method) => <option key={method}>{method}</option>)}</select></label>
      <Input name="endpoint" label={__('Restricted endpoint')} required placeholder="/identity/restricted/accounts/${accountId}" dataTestId="template-endpoint"/>
      <label>{__('Scope')}<select name="scope" data-testid="template-scope"><option>GLOBAL</option><option>SCHOOL</option></select></label>
      <label>{__('Headers (JSON object)')}<textarea name="headers" defaultValue="{}" required rows={3} data-testid="template-headers"/></label>
      <label>{__('Input schema (JSON object)')}<textarea name="schema" defaultValue="{}" required rows={6} data-testid="template-schema"/></label>
    </ActionForm>}
    {environments && <PagedRecords<Environment> resource="environments" title={__('Environments')} revision={revision} rowKey={(r) => r.id} columns={[{title:__('Code'),cell:(r)=>r.code},{title:__('Name'),cell:(r)=>r.name},{title:__('ID'),cell:(r)=>r.id}]}/>}
    {!environments && <PagedRecords<Template> resource="api-templates" title={__('Template revisions')} revision={revision} rowKey={(r)=>r.id} onSelect={setSelected} columns={[{title:__('Code'),cell:(r)=>r.code},{title:__('Version'),cell:(r)=>r.version},{title:__('Method'),cell:(r)=>r.method},{title:__('Endpoint'),cell:(r)=>r.endpoint},{title:__('Scope'),cell:(r)=>r.scope}]}/>}
    {selected && !environments && <section className="action-card"><div className="page-heading"><h2>{selected.code} · v{selected.version}</h2><Button onClick={()=>setSelected(null)} dataTestId="template-detail-close">{__('Close')}</Button></div><pre>{JSON.stringify(selected,null,2)}</pre></section>}
  </>;
}
