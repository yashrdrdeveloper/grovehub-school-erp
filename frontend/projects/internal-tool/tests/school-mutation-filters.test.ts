import { test } from 'node:test';
import assert from 'node:assert/strict';
import { buildReceiptFilters } from '../packages/domain/audit/pages/helpers/schoolMutation.filters.ts';
const id='e93696cc-15f0-4e06-9577-0d5cfc2a83c8';
test('school mutation filters preserve the selected school and operation',()=>{
 assert.deepEqual(buildReceiptFilters(id,'DISABLE'),{schoolId:id,operation:'DISABLE'});
 assert.deepEqual(buildReceiptFilters(id,''),{schoolId:id});
 assert.deepEqual(buildReceiptFilters('','RENAME'),{operation:'RENAME'});
 assert.deepEqual(buildReceiptFilters('',''),{});
});
test('invalid route filters never broaden the search',()=>{
 for(const value of ['invalid',id.toUpperCase(),'00000000-0000-0000-0000-000000000000']){
  assert.equal(buildReceiptFilters(value,''),null);
 }
 assert.equal(buildReceiptFilters(id,'DELETE'),null);
});
