import { useState } from 'react';
import { Select } from '@Components/components/Select';
import { TextAreaControl } from '@Components/components/TextAreaControl';
import { ApiError, requestJSON } from '../../../../business/api/client';
import { ActionForm } from '../../../../shell/components/ActionForm';
import { SchemaField } from './components/SchemaField';
import { serializeRequest } from './helpers/requestForm.serialize';
import { getRequestTemplates, getRequestEnvironments, getRequestSchoolIds } from './helpers/requestForm.grants';
import { useRequestRegistry } from './hooks/useRequestRegistry';
import type { RequestFormProps } from './types/requestForm.props';

export function RequestForm({ close, success }: RequestFormProps) {
  const { templates, environments, grants, error, ready } = useRequestRegistry();
  const [templateId, setTemplateId] = useState('');
  const [environmentId, setEnvironmentId] = useState('');
  const [schoolId, setSchoolId] = useState('');
  const availableTemplates = getRequestTemplates(templates, grants);
  const template = availableTemplates.find((item) => item.id === templateId);
  const availableEnvironments = getRequestEnvironments(template, environments, grants);
  const environment = availableEnvironments.find((item) => item.id === environmentId);
  const schools = getRequestSchoolIds(template, environment, grants);
  const schoolScoped = template?.scope === 'SCHOOL';
  const rename = schoolScoped && template?.method === 'PATCH' && template.endpoint === '/school/restricted/registry/schools/${schoolId}/name';
  const selectedGrant = grants.find((grant) => grant.canRequest && grant.schoolId === schoolId && grant.templateId === templateId && grant.environmentCode === environment?.code);
  const hasSnapshot = Boolean(selectedGrant?.schoolName && selectedGrant?.schoolVersion && /^[1-9][0-9]*$/.test(selectedGrant.schoolVersion));
  const targetReady = Boolean(ready && template && environment && (!schoolScoped || schools.includes(schoolId)) && (!rename || hasSnapshot));

  function changeTemplate(value: string | null) {
    setTemplateId(value ?? '');
    setEnvironmentId('');
    setSchoolId('');
  }
  function changeEnvironment(value: string | null) {
    setEnvironmentId(value ?? '');
    setSchoolId('');
  }
  async function submit(data: FormData) {
    if (!targetReady || !template || !environment) throw new ApiError(400, 'INVALID_REQUEST');
    data.set('environmentId', environment.id);
    if (schoolScoped) {
      data.set('schoolId', schoolId);
      // The pilot's schoolId placeholder comes from the same authorized target.
      data.set('mode:schoolId', 'value');
      data.set('value:schoolId', schoolId);
    }
    if (rename) {
      if (!selectedGrant?.schoolName || !selectedGrant.schoolVersion) throw new ApiError(400, 'INVALID_REQUEST');
      data.set('mode:expectedName', 'value');
      data.set('value:expectedName', selectedGrant.schoolName);
      data.set('mode:expectedVersion', 'value');
      data.set('value:expectedVersion', selectedGrant.schoolVersion);
    }
    const raw = serializeRequest(template.id, schoolScoped, template.bodySchema, data);
    return requestJSON('/internal/protected/requests', 'POST', raw);
  }

  return <ActionForm title={__('Submit approval request')} onClose={close} onSuccess={success} submit={submit} submitDisabled={!targetReady}>
    {error && <p role="alert" className="error">{error}</p>}
    {!ready && !error && <p role="status">{__('Loading authorized request targets…')}</p>}
    {ready && !availableTemplates.length && <p role="status">{__('No request templates available for your permissions.')}</p>}
    <Select label={__('Template revision')} value={templateId || null}
      options={availableTemplates.map((item) => ({ value: item.id, label: item.code + ' · v' + item.version }))}
      disabled={!ready} searchable clearable={false} onChange={changeTemplate} dataTestId="request-template" />
    <Select label={__('Environment')} value={environmentId || null}
      options={availableEnvironments.map((item) => ({ value: item.id, label: item.name }))}
      disabled={!template || !availableEnvironments.length} searchable clearable={false} onChange={changeEnvironment} dataTestId="request-environment" />
    {schoolScoped && <Select label={__('Authorized school')} value={schoolId || null}
      options={schools.map((id) => {
        const grant = grants.find((item) => item.schoolId === id && item.templateId === templateId && item.environmentCode === environment?.code && item.canRequest);
        return { value: id, label: grant?.schoolName ? grant.schoolName + ' · ' + grant.schoolCode : id };
      })}
      disabled={!environment || !schools.length} searchable clearable={false}
      onChange={(value) => setSchoolId(value ?? '')} dataTestId="request-school" />}
    {schoolScoped && environment && !schools.length && <p role="status">{__('No authorized school for this template and environment.')}</p>}
    {rename && selectedGrant && hasSnapshot && <section data-testid="request-school-snapshot">
      <p>{__('Current school name')}: {selectedGrant.schoolName}</p>
      <p className="muted">{__('School version')}: {selectedGrant.schoolVersion}</p>
      <p className="muted">{__('The approved request is valid only for this exact school version.')}</p>
    </section>}
    {rename && schoolId && !hasSnapshot && <p role="alert">{__('Current school details are unavailable. Reopen the request form.')}</p>}
    {template && <details><summary>{__('Input schema')}</summary><pre>{JSON.stringify(template.bodySchema, null, 2)}</pre></details>}
    {targetReady && template && Object.entries(template.bodySchema).filter(([name]) => (!schoolScoped || name !== 'schoolId') && (!rename || (name !== 'expectedName' && name !== 'expectedVersion'))).map(([name, field]) =>
      <SchemaField key={[template.id, environmentId, schoolId, name].join(':')} name={name} field={field} />)}
    <label>{__('Reason')}<TextAreaControl name="reason" required maxLength={1000} rows={3} dataTestId="request-reason" /></label>
  </ActionForm>;
}
