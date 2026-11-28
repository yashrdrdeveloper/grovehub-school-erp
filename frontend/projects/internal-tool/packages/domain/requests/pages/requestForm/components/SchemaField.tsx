import { useState } from 'react';
import { Input } from '@Components/components/Input';
import { Select } from '@Components/components/Select';
import { TextAreaControl } from '@Components/components/TextAreaControl';
import type { SchemaFieldProps } from '../types/requestForm.props';

function initialMode(required: boolean) {
  if (required) return 'value';
  return 'omit';
}

export function SchemaField({ name, field }: SchemaFieldProps) {
  const [mode, setMode] = useState(initialMode(field.required));
  const [booleanValue, setBooleanValue] = useState('true');
  const modes = [
    { value: 'value', label: __('Provide value') },
    ...(!field.required && [{ value: 'omit', label: __('Omit field') }] || []),
    ...(field.nullable && [{ value: 'null', label: __('Null') }] || []),
  ];
  const testId = 'schema-field-' + name;
  return (
    <fieldset data-testid={testId}>
      <legend>{name} · {field.type}{field.required && ' *'}</legend>
      {modes.length > 1 && <Select label={__('Field value mode')} options={modes} value={mode} clearable={false} searchable={false} onChange={(value) => { if (value) setMode(value); }} dataTestId={testId + '-mode'} />}
      <input type="hidden" name={'mode:' + name} value={mode} />
      {mode === 'value' && field.type === 'string' && <Input name={'value:' + name} label={name} maxLength={4096} dataTestId={testId + '-value'} />}
      {mode === 'value' && (field.type === 'number' || field.type === 'integer') && <Input name={'value:' + name} label={name} inputMode="decimal" required dataTestId={testId + '-value'} />}
      {mode === 'value' && field.type === 'boolean' && <>
        <Select label={name} value={booleanValue} options={[{value:'true',label:__('True')},{value:'false',label:__('False')}]} clearable={false} searchable={false} onChange={(value) => { if (value) setBooleanValue(value); }} dataTestId={testId + '-boolean'} />
        <input type="hidden" name={'value:' + name} value={booleanValue} />
      </>}
      {mode === 'value' && (field.type === 'object' || field.type === 'array') && <label>
        {__('JSON value')}
        <TextAreaControl name={'value:' + name} required rows={4} defaultValue={(field.type === 'array' && '[]') || '{}'} dataTestId={testId + '-json'} />
      </label>}
    </fieldset>
  );
}
