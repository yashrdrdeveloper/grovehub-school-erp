# Platform API migration inventory

Generated from existing controller annotations. Paths below are controller-relative to the service and tier; verify actual route registrations before cutover. All target Platform endpoints use restricted workload authentication. Wrapper paths are proposed, not implemented. Preserve browser contracts through application-owned compatibility routes where necessary.

| Method | Existing controller route | Wrapper owner | Proposed Platform route |
|---|---|---|---|
| POST | /audit/protected/search | School | /audit/restricted/search |
| POST | /foundation/protected/persons | School | /foundation/restricted/persons |
| GET | /foundation/protected/persons/{gid} | School | /foundation/restricted/persons/{gid} |
| PATCH | /foundation/protected/persons/{gid} | School | /foundation/restricted/persons/{gid} |
| POST | /foundation/protected/persons/search | School | /foundation/restricted/persons/search |
| PATCH | /foundation/protected/persons/{gid}/deactivate | School | /foundation/restricted/persons/{gid}/deactivate |
| POST | /identity/protected/accounts | Identity | /identity/restricted/accounts |
| POST | /identity/protected/accounts/search | Identity | /identity/restricted/accounts/search |
| GET | /identity/protected/accounts/roles | Identity | /identity/restricted/accounts/roles |
| GET | /identity/protected/accounts/{accountId} | Identity | /identity/restricted/accounts/{accountId} |
| PATCH | /identity/protected/accounts/{accountId} | Identity | /identity/restricted/accounts/{accountId} |
| PATCH | /identity/protected/accounts/{accountId}/suspend | Identity | /identity/restricted/accounts/{accountId}/suspend |
| PATCH | /identity/protected/accounts/{accountId}/activate | Identity | /identity/restricted/accounts/{accountId}/activate |
| PATCH | /identity/protected/accounts/{accountId}/deactivate | Identity | /identity/restricted/accounts/{accountId}/deactivate |
| POST | /identity/protected/accounts/{accountId}/resend-invitation | Identity | /identity/restricted/accounts/{accountId}/resend-invitation |
| GET | /identity/protected/account-roles/{accountId} | Identity | /identity/restricted/account-roles/{accountId} |
| POST | /identity/protected/account-roles/{accountId}/preview | Identity | /identity/restricted/account-roles/{accountId}/preview |
| GET | /identity/public/accounts/setup/status | Identity | /identity/restricted/accounts/setup/status |
| POST | /identity/public/accounts/setup | Identity | /identity/restricted/accounts/setup |
| POST | /identity/protected/approvals/history/search | Identity | /identity/restricted/approvals/history/search |
| POST | /identity/protected/approvals/inbox/search | Identity | /identity/restricted/approvals/inbox/search |
| GET | /identity/protected/approvals/inbox/count | Identity | /identity/restricted/approvals/inbox/count |
| GET | /identity/protected/approvals/{requestId} | Identity | /identity/restricted/approvals/{requestId} |
| POST | /identity/protected/approvals/lifecycle | Identity | /identity/restricted/approvals/lifecycle |
| POST | /identity/protected/approvals/permissions | Identity | /identity/restricted/approvals/permissions |
| POST | /identity/protected/approvals/assignments | Identity | /identity/restricted/approvals/assignments |
| POST | /identity/protected/approvals/{requestId}/approve | Identity | /identity/restricted/approvals/{requestId}/approve |
| POST | /identity/protected/approvals/{requestId}/reject | Identity | /identity/restricted/approvals/{requestId}/reject |
| POST | /identity/protected/approvals/{requestId}/cancel | Identity | /identity/restricted/approvals/{requestId}/cancel |
| POST | /identity/protected/approvals/{requestId}/execute | Identity | /identity/restricted/approvals/{requestId}/execute |
| POST | /identity/protected/auth/logout | Identity | /identity/restricted/auth/logout |
| POST | /identity/protected/auth/logout-all | Identity | /identity/restricted/auth/logout-all |
| POST | /identity/protected/auth/change-password | Identity | /identity/restricted/auth/change-password |
| POST | /identity/protected/auth/reauthenticate | Identity | /identity/restricted/auth/reauthenticate |
| GET | /identity/protected/auth/sessions | Identity | /identity/restricted/auth/sessions |
| DELETE | /identity/protected/auth/sessions/{sessionId} | Identity | /identity/restricted/auth/sessions/{sessionId} |
| GET | /identity/protected/auth/me | Identity | /identity/restricted/auth/me |
| POST | /identity/public/auth/login | Identity | /identity/restricted/auth/login |
| POST | /identity/public/auth/refresh | Identity | /identity/restricted/auth/refresh |
| POST | /identity/public/auth/forgot-password | Identity | /identity/restricted/auth/forgot-password |
| GET | /identity/public/auth/recovery/validate | Identity | /identity/restricted/auth/recovery/validate |
| GET | /identity/public/auth/invitations/validate | Identity | /identity/restricted/auth/invitations/validate |
| POST | /identity/public/auth/recovery/exchange | Identity | /identity/restricted/auth/recovery/exchange |
| GET | /identity/public/auth/reset-password/status | Identity | /identity/restricted/auth/reset-password/status |
| POST | /identity/public/auth/reset-password | Identity | /identity/restricted/auth/reset-password |
| POST | /identity/public/auth/invitations/exchange | Identity | /identity/restricted/auth/invitations/exchange |
| GET | /identity/protected/capabilities/me | Identity | /identity/restricted/capabilities/me |
| GET | /identity/protected/capabilities/permissions | Identity | /identity/restricted/capabilities/permissions |
| POST | /identity/protected/roles | Identity | /identity/restricted/roles |
| PATCH | /identity/protected/roles/{roleId} | Identity | /identity/restricted/roles/{roleId} |
| POST | /identity/protected/roles/search | Identity | /identity/restricted/roles/search |
| GET | /identity/protected/roles/{roleId} | Identity | /identity/restricted/roles/{roleId} |
| GET | /identity/protected/roles/{roleId}/permissions | Identity | /identity/restricted/roles/{roleId}/permissions |
| POST | /identity/protected/roles/{roleId}/permissions/preview | Identity | /identity/restricted/roles/{roleId}/permissions/preview |
| POST | /integration/protected/connections | School | /integration/restricted/connections |
| GET | /integration/protected/connections/{id} | School | /integration/restricted/connections/{id} |
| GET | /integration/protected/connections | School | /integration/restricted/connections |
| DELETE | /integration/protected/connections/{id} | School | /integration/restricted/connections/{id} |
| POST | /integration/protected/connections/{id}/oauth/start | School | /integration/restricted/connections/{id}/oauth/start |
| POST | /integration/protected/connections/{id}/test | School | /integration/restricted/connections/{id}/test |
| GET | /integration/protected/connections/{id}/tests | School | /integration/restricted/connections/{id}/tests |
| POST | /integration/protected/connections/{id}/watch | School | /integration/restricted/connections/{id}/watch |
| POST | /integration/protected/connections/system-mail/sender | School | /integration/restricted/connections/system-mail/sender |
| POST | /integration/protected/connections/{id}/actions/{action} | School | /integration/restricted/connections/{id}/actions/{action} |
| GET | /integration/public/connections/oauth/{provider}/callback | Verified ingress | /integration/restricted/connections/oauth/{provider}/callback |
| GET | /integration/public/providers | School | /integration/restricted/providers |
| GET | /integration/public/providers/{provider} | School | /integration/restricted/providers/{provider} |

Total: 67 annotated operations. Existing restricted routes require separate collision and policy review. OAuth callback must terminate at verified ingress; do not blindly forward browser credentials. Public Identity wrappers retain rate limits and recovery protections. Protected wrappers enforce session, CSRF, tenant and capability policies. Platform verifies workload and delegated actor independently.
