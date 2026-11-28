import { test } from 'node:test';
import assert from 'node:assert/strict';
import { getRequestTemplates, getRequestEnvironments, getRequestSchoolIds } from '../packages/domain/requests/pages/requestForm/helpers/requestForm.grants.ts';
import type { SchoolGrant, Template, Environment } from '../packages/business/api/types.ts';

const school: Template = { id:'revision1', code:'check', version:1, module:'audit', method:'GET', endpoint:'/audit/restricted/schools/${schoolId}/activity-check', scope:'SCHOOL',headers:{},bodySchema:{} };
const global: Template = {...school,id:'global',scope:'GLOBAL'};
const next: Template = {...school,id:'revision2',version:2};
const dev: Environment = {id:'env1',code:'dev',name:'Development'};
const prod: Environment = {id:'env2',code:'prod',name:'Production'};
const allowed: SchoolGrant = {environmentCode:'dev',templateId:school.id,schoolId:'school-a',canRequest:true,canReview:false,canExecute:false};
const grants = [allowed,{...allowed,schoolId:'school-b',canRequest:false,canReview:true},{...allowed,schoolId:'school-c'},{...allowed,schoolId:'school-d',environmentCode:'prod',templateId:next.id}];

test('only requester grants unlock exact school template revisions and environments',()=>{
  assert.deepEqual(getRequestTemplates([school,next,global],grants),[school,next,global]);
  assert.deepEqual(getRequestTemplates([school,next,global],[]),[global]);
  assert.deepEqual(getRequestEnvironments(school,[dev,prod],grants),[dev]);
  assert.deepEqual(getRequestSchoolIds(school,dev,grants),['school-a','school-c']);
  assert.deepEqual(getRequestSchoolIds(school,prod,grants),[]);
  assert.deepEqual(getRequestSchoolIds(next,prod,grants),['school-d']);
  assert.deepEqual(getRequestSchoolIds(global,dev,grants),[]);
  assert.deepEqual(getRequestEnvironments(global,[dev,prod],grants),[dev,prod]);
});
