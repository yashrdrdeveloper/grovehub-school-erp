import { test } from 'node:test';
import assert from 'node:assert/strict';
import { csrfFromCookie } from '../packages/business/api/client.ts';
import { parseObject } from '../packages/business/api/helpers.ts';

test('only Internal Tool CSRF cookie is accepted',()=>{
 assert.equal(csrfFromCookie('gh_csrf=school; __Host-gh_internal_csrf=internal'),'internal');
 assert.equal(csrfFromCookie('gh_csrf=school'),'');
 assert.equal(csrfFromCookie('__Host-gh_internal_csrf=a; __Host-gh_internal_csrf=b'),'');
});
test('input must be an object',()=>{
 assert.deepEqual(parseObject('{"enabled":true}'),{enabled:true});
 for(const input of ['[]','null','1','invalid'])assert.throws(()=>parseObject(input));
});
