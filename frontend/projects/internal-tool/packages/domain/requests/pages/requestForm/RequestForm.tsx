import { useState } from 'react';
import { Input } from '@Components/components/Input';
import { Select } from '@Components/components/Select';
import { TextAreaControl } from '@Components/components/TextAreaControl';
import { ApiError, requestJSON } from '../../../../business/api/client';
import { ActionForm } from '../../../../shell/components/ActionForm';
import { SchemaField } from './components/SchemaField';
import { serializeRequest } from './helpers/requestForm.serialize';
import { useRequestRegistry } from './hooks/useRequestRegistry';
import type { RequestFormProps } from './types/requestForm.props';

export function RequestForm({ close, success }: RequestFormProps) {
  const { templates, environments, error, ready } = useRequestRegistry();
  const [templateId, setTemplateId] = useState('');
  const [environmentId, setEnvironmentId] = useState('');
  const template = templates.find((item) => item.id === templateId);

  async function submit(data: FormData) {
    if (!ready || !template || !environmentId) throw new ApiError(400, 'INVALID_REQUEST');
    const raw = serializeRequest(template.id, template.scope === 'SCHOOL', template.bodySchema, data);
    return requestJSON('/internal/protected/requests', 'POST', raw);
  }

  return <ActionForm title={__('Submit approval request')} onClose={close} onSuccess={success} submit={submit}>
    {error && <p role="alert" className="error">{error}</p>}
    <Select label={__('Template revision')} value={templateId || null}
      options={templates.map((item) => ({ value: item.id, label: item.code + ' · v' + item.version }))}
      disabled={!ready} searchable clearable={false} onChange={(value) => setTemplateId(value ?? '')} dataTestId="request-template" />
    <Select label={__('Environment')} value={environmentId || null}
      options={environments.map((item) => ({ value: item.id, label: item.name }))}
      disabled={!ready} searchable clearable={false} onChange={(value) => setEnvironmentId(value ?? '')} dataTestId="request-environment" />
    <input type="hidden" name="environmentId" value={environmentId} />
    {template?.scope === 'SCHOOL' && <Input name="schoolId" label={__('School ID')} required dataTestId="request-school" />}
    {template && <details><summary>{__('Input schema')}</summary><pre>{JSON.stringify(template.bodySchema, null, 2)}</pre></details>}
    {template && Object.entries(template.bodySchema).map(([name, field]) => <SchemaField key={template.id + ':' + name} name={name} field={field} />)}
    <label>{__('Reason')}<TextAreaControl name="reason" required maxLength={1000} rows={3} dataTestId="request-reason" /></label>
  </ActionForm>;
}
