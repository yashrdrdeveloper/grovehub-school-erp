import { test } from 'node:test';
import assert from 'node:assert/strict';
import { serializeSchemaInput, serializeRequest } from '../packages/domain/requests/pages/requestForm/helpers/requestForm.serialize.ts';

test('preserve exact int64 and nested numeric JSON without JS rounding', () => {
  const data = new FormData();
  data.set('mode:version','value'); data.set('value:version','9223372036854775807');
  data.set('mode:nested','value'); data.set('value:nested','{"value":9007199254740993}');
  data.set('mode:optional','omit');
  data.set('mode:nullable','null');
  const schema = { version: {type:'integer',required:true}, nested:{type:'object',required:true}, optional:{type:'string',required:false}, nullable:{type:'string',required:true,nullable:true} };
  assert.equal(serializeSchemaInput(schema,data),'{"version":9223372036854775807,"nested":{"value":9007199254740993},"nullable":null}');
  assert.match(serializeRequest('template',false,schema,data),/"input":\{"version":9223372036854775807/);
  data.set('value:version','9223372036854775808');
  assert.throws(() => serializeSchemaInput(schema,data));
});
test('required and non-nullable fields cannot be omitted or null', () => {
  const data = new FormData();
  const schema = { id:{type:'string',required:true} };
  for (const mode of ['omit','null']) { data.set('mode:id',mode); assert.throws(()=>serializeSchemaInput(schema,data)); }
});
