import { test } from 'node:test';
import assert from 'node:assert/strict';
import { buildSchoolRecords, filterSchools } from '../packages/domain/schools/pages/helpers/schools.records.ts';
import type { Environment, SchoolGrant, Template } from '../packages/business/api/types.ts';
const env: Environment = {id:'dev',code:'DEV',name:'Development'};
const prod: Environment = {id:'prod',code:'PROD',name:'Production'};
const rename: Template = {id:'rename',code:'rename',version:1,module:'school',method:'PATCH',endpoint:'/school/restricted/registry/schools/${schoolId}/name',scope:'SCHOOL',headers:{},bodySchema:{}};
const lifecycle: Template = {...rename,id:'lifecycle',code:'status',endpoint:'/school/restricted/registry/schools/${schoolId}/status'};
const grant: SchoolGrant = {environmentCode:'DEV',templateId:'rename',schoolId:'school',schoolName:'Example School',schoolCode:'EX',schoolVersion:'9007199254740993',schoolStatus:'ACTIVE',canRequest:true,canReview:false,canExecute:false};
test('school records preserve version precision and exact environment and template grants',()=>{
 const records=buildSchoolRecords([grant,{...grant,templateId:'lifecycle',lifecycle:true},{...grant,environmentCode:'PROD',canRequest:false,canReview:true}],[env,prod],[rename,lifecycle]);
 assert.equal(records.length,2);
 assert.equal(records[0]?.version,'9007199254740993');
 assert.deepEqual(records[0]?.actions.map((item)=>item.kind),['rename','lifecycle']);
 assert.equal(records[1]?.actions.length,0);
 assert.equal(filterSchools(records,'ex','prod','ACTIVE').length,1);
 assert.equal(filterSchools(records,'missing','','').length,0);
});
test('disabled schools only expose explicitly granted lifecycle requests',()=>{
 const disabled: SchoolGrant={...grant,schoolStatus:'DISABLED'};
 const records=buildSchoolRecords([disabled,{...disabled,templateId:'lifecycle',lifecycle:true}],[env],[rename,lifecycle]);
 assert.deepEqual(records[0]?.actions.map((item)=>item.kind),['lifecycle']);
});
test('unresolved and inconsistent school snapshots cannot request changes',()=>{
 for(const grants of [[{...grant,schoolVersion:undefined}],[grant,{...grant,templateId:'lifecycle',schoolVersion:'2',lifecycle:true}]]){
  const records=buildSchoolRecords(grants,[env],[rename,lifecycle]);
  assert.equal(records[0]?.snapshotReady,false);
  assert.deepEqual(records[0]?.actions,[]);
 }
});
